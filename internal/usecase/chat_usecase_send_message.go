package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/pkg/logger"
)

func (c *ChatUseCase) SendMessage(ctx context.Context, userId int, sessionId int64, userMessage string, attachmentFileIDs []int64, fileRAG *SendMessageFileRAGOptions) (chan ChatStreamChunk, error) {
	logger.I("SendMessage: phase=enter session_id=%d user_id=%d attachments=%d file_rag=%t", sessionId, userId, len(attachmentFileIDs), fileRAG != nil)
	session, err := c.verifySessionOwnership(ctx, userId, sessionId)
	if err != nil {
		logger.W("SendMessage: сессия не принадлежит пользователю: %v", err)
		return nil, err
	}

	runnerAddr, resolvedModel, err := c.chatRunnerAddrAndModel(ctx, session)
	if err != nil {
		logger.W("SendMessage: раннер/модель: %v", err)
		return nil, err
	}

	messages, err := c.historyMessagesForLLM(ctx, sessionId)
	if err != nil {
		logger.E("SendMessage: история для LLM: %v", err)
		return nil, err
	}

	normalizedAttachmentFileIDs, err := normalizeAttachmentFileIDsForModel(attachmentFileIDs)
	if err != nil {
		return nil, err
	}

	if err := validateFileRAGOptions(fileRAG, normalizedAttachmentFileIDs); err != nil {
		return nil, err
	}

	if len(normalizedAttachmentFileIDs) == 0 && strings.TrimSpace(userMessage) == "" {
		return nil, fmt.Errorf("пустое сообщение: укажите текст или вложение")
	}

	attachmentNames := make([]string, 0, len(normalizedAttachmentFileIDs))
	attachmentContents := make([][]byte, 0, len(normalizedAttachmentFileIDs))
	attachmentImageMIMEs := make([]string, 0, len(normalizedAttachmentFileIDs))
	for _, fid := range normalizedAttachmentFileIDs {
		name, content, imgMime, err := c.loadSessionAttachmentForSend(ctx, userId, sessionId, fid)
		if err != nil {
			return nil, err
		}
		attachmentNames = append(attachmentNames, name)
		attachmentContents = append(attachmentContents, content)
		attachmentImageMIMEs = append(attachmentImageMIMEs, imgMime)
	}

	var storedAttachmentFileID *int64
	if len(normalizedAttachmentFileIDs) > 0 {
		v := normalizedAttachmentFileIDs[0]
		storedAttachmentFileID = &v
	}
	userMsg := domain.NewMessageWithAttachment(sessionId, userMessage, domain.MessageRoleUser, storedAttachmentFileID)
	if err := c.messageRepo.Create(ctx, userMsg); err != nil {
		logger.E("SendMessage: создание сообщения: %v", err)
		return nil, err
	}

	settings, _ := c.sessionSettingsRepo.GetBySessionID(ctx, sessionId)
	stopSequences, timeoutSeconds, genParams := genParamsFromSessionSettings(settings)
	c.injectWebSearchAndMCP(ctx, genParams, settings, userId, sessionId)

	messagesForLLM, ragStream, err := c.buildSendPromptAssembly(
		ctx,
		sendPromptAssemblyInput{
			sessionID:                sessionId,
			userID:                   userId,
			resolvedModel:            resolvedModel,
			settings:                 settings,
			history:                  messages,
			userMessage:              userMessage,
			userMsg:                  userMsg,
			attachmentFileIDs:        normalizedAttachmentFileIDs,
			attachmentNames:          attachmentNames,
			attachmentContents:       attachmentContents,
			attachmentImageMIMEs:     attachmentImageMIMEs,
			fileRAG:                  fileRAG,
			preferFullDocumentIfFits: c.preferFullDocumentWhenFits,
			genParams:                genParams,
		},
	)
	if err != nil {
		return nil, err
	}

	if err := c.hydrateAttachmentsForRunner(ctx, messagesForLLM); err != nil {
		logger.E("SendMessage: подгрузка вложений для раннера: %v", err)
		return nil, err
	}
	var historyNotice bool
	messagesForLLM, historyNotice = c.capLLMHistoryTokens(ctx, messagesForLLM, 1, sessionId, resolvedModel, runnerAddr, true)

	if genParams != nil && len(genParams.Tools) > 0 {
		logger.I("SendMessage: phase=branch_tool_loop session_id=%d user_id=%d runner=%q model=%q tools=%d history_notice=%t", sessionId, userId, runnerAddr, resolvedModel, len(genParams.Tools), historyNotice)
		return c.sendMessageWithToolLoop(ctx, userId, sessionId, runnerAddr, resolvedModel, messagesForLLM, stopSequences, timeoutSeconds, genParams, historyNotice, ragStream)
	}
	if settings != nil {
		logger.W("SendMessage: tool-loop disabled session_id=%d tools=%d mcp_enabled=%t mcp_server_ids=%v web_search_enabled=%t",
			sessionId,
			func() int {
				if genParams == nil {
					return 0
				}

				return len(genParams.Tools)
			}(),
			settings.MCPEnabled,
			settings.MCPServerIDs,
			settings.WebSearchEnabled,
		)
	} else {
		logger.W("SendMessage: tool-loop disabled session_id=%d settings=nil tools=0", sessionId)
	}

	assistantMsg := domain.NewMessage(sessionId, "", domain.MessageRoleAssistant)
	if err := c.messageRepo.Create(ctx, assistantMsg); err != nil {
		logger.E("SendMessage: создание черновика ответа: %v", err)
		return nil, err
	}
	messageID := assistantMsg.Id

	logger.I("SendMessage: phase=branch_plain_llm_stream session_id=%d user_id=%d runner=%q model=%q tools=%d",
		sessionId, userId, runnerAddr, resolvedModel, func() int {
			if genParams == nil {
				return 0
			}
			return len(genParams.Tools)
		}())
	responseChan, err := c.llmRepo.SendMessageOnRunner(ctx, runnerAddr, sessionId, resolvedModel, messagesForLLM, stopSequences, timeoutSeconds, genParams)
	if err != nil {
		logger.E("SendMessage: вызов LLM: %v", err)
		return nil, err
	}
	logger.I("SendMessage: phase=plain_llm_stream_started session_id=%d message_id=%d", sessionId, messageID)

	var fullResponse strings.Builder
	clientChan := make(chan ChatStreamChunk, 100)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.E("SendMessage panic: session_id=%d user_id=%d panic=%v", sessionId, userId, r)
				select {
				case <-ctx.Done():
				case clientChan <- ChatStreamChunk{Kind: StreamChunkKindText, Text: "внутренняя ошибка обработки ответа"}:
				}
			}
		}()

		defer func() {
			_ = c.messageRepo.UpdateContent(context.Background(), messageID, fullResponse.String())
		}()

		defer close(clientChan)

		if ragStream != nil {
			select {
			case <-ctx.Done():
				return
			case clientChan <- ragStream.asChunk():
			}
		}

		if historyNotice {
			select {
			case <-ctx.Done():
				return
			case clientChan <- ChatStreamChunk{Kind: StreamChunkKindNotice, Text: HistoryTruncatedClientNotice}:
			}
		}

		forwardLLMStreamChunks(ctx, clientChan, messageID, responseChan, &fullResponse)
	}()

	return clientChan, nil
}
