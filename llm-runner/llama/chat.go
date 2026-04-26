package llama

/*
#include "wrapper.h"
#include <stdlib.h>
*/
import "C"

import (
	gocontext "context"
	"fmt"
	"strings"
	"unsafe"
)

func formatChatMessages(model *Model, messages []ChatMessage, opts ChatOptions) (string, error) {
	template := opts.ChatTemplate
	if template == "" {
		template = model.ChatTemplate()
	}
	if template == "" {
		return "", fmt.Errorf("chat-шаблон недоступен: укажите ChatOptions.ChatTemplate или используйте модель со встроенным шаблоном (или Generate() для обычной генерации)")
	}

	prompt, err := applyChatTemplate(template, messages, true)
	if err != nil {
		return "", fmt.Errorf("не удалось применить chat-шаблон: %w", err)
	}

	return prompt, nil
}

func parseReasoning(text string, format ReasoningFormat, chatFormat int) (content, reasoningContent string, err error) {
	if format == ReasoningFormatNone || text == "" {
		return text, "", nil
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cFormat := C.llama_wrapper_reasoning_format(format)
	cChatFormat := C.int(chatFormat)

	result := C.llama_wrapper_parse_reasoning(cText, C.bool(true), cFormat, cChatFormat)
	if result == nil {
		return "", "", fmt.Errorf("не удалось распарсить reasoning: %s", C.GoString(C.llama_wrapper_last_error()))
	}
	defer C.llama_wrapper_free_parsed_message(result)

	content = C.GoString(result.content)
	if result.reasoning_content != nil {
		reasoningContent = C.GoString(result.reasoning_content)
	}

	return content, reasoningContent, nil
}

func (m *Model) chatWithContext(ctx gocontext.Context, c *Context, messages []ChatMessage, opts ChatOptions) (*ChatResponse, error) {
	prompt, err := formatChatMessages(m, messages, opts)
	if err != nil {
		return nil, err
	}

	genOpts := []GenerateOption{
		WithStopWords(opts.StopWords...),
	}

	if opts.MaxTokens != nil {
		genOpts = append(genOpts, WithMaxTokens(*opts.MaxTokens))
	}
	if opts.Temperature != nil {
		genOpts = append(genOpts, WithTemperature(*opts.Temperature))
	}

	if opts.TopP != nil {
		genOpts = append(genOpts, WithTopP(*opts.TopP))
	}

	if opts.TopK != nil {
		genOpts = append(genOpts, WithTopK(*opts.TopK))
	}

	if opts.Seed != nil {
		genOpts = append(genOpts, WithSeed(*opts.Seed))
	}

	tokenCh, errCh := c.GenerateChannel(ctx, prompt, genOpts...)

	var content strings.Builder

Loop:
	for {
		select {
		case token, ok := <-tokenCh:
			if !ok {
				break Loop
			}
			content.WriteString(token)
		case err := <-errCh:
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	fullOutput := content.String()
	chatFormat := m.getChatFormat()
	parsedContent, reasoning, err := parseReasoning(fullOutput, opts.ReasoningFormat, chatFormat)
	if err != nil {
		return &ChatResponse{Content: fullOutput}, nil
	}

	return &ChatResponse{
		Content:          parsedContent,
		ReasoningContent: reasoning,
	}, nil
}

func (m *Model) chatStreamWithContext(ctx gocontext.Context, c *Context, messages []ChatMessage, opts ChatOptions) (<-chan ChatDelta, <-chan error) {
	bufferSize := 256
	if opts.StreamBufferSize > 0 {
		bufferSize = opts.StreamBufferSize
	}

	deltaCh := make(chan ChatDelta, bufferSize)
	errCh := make(chan error, 1)

	go func() {
		defer close(deltaCh)
		defer close(errCh)

		prompt, err := formatChatMessages(m, messages, opts)
		if err != nil {
			select {
			case errCh <- err:
			default:
			}
			return
		}

		genOpts := []GenerateOption{
			WithStopWords(opts.StopWords...),
		}

		if opts.MaxTokens != nil {
			genOpts = append(genOpts, WithMaxTokens(*opts.MaxTokens))
		}

		if opts.Temperature != nil {
			genOpts = append(genOpts, WithTemperature(*opts.Temperature))
		}

		if opts.TopP != nil {
			genOpts = append(genOpts, WithTopP(*opts.TopP))
		}

		if opts.TopK != nil {
			genOpts = append(genOpts, WithTopK(*opts.TopK))
		}

		if opts.Seed != nil {
			genOpts = append(genOpts, WithSeed(*opts.Seed))
		}

		tokenCh, genErrCh := c.GenerateChannel(ctx, prompt, genOpts...)

		chatFormat := m.getChatFormat()

		var accumulated strings.Builder
		var prevContent, prevReasoning string

	Loop:
		for {
			select {
			case token, ok := <-tokenCh:
				if !ok {
					break Loop
				}

				accumulated.WriteString(token)

				content, reasoning, err := parseReasoning(accumulated.String(), opts.ReasoningFormat, chatFormat)
				if err != nil {
					select {
					case deltaCh <- ChatDelta{Content: token}:
					case <-ctx.Done():
						return
					}
					continue
				}

				contentDelta := content[len(prevContent):]
				reasoningDelta := reasoning[len(prevReasoning):]

				if contentDelta != "" || reasoningDelta != "" {
					select {
					case deltaCh <- ChatDelta{
						Content:          contentDelta,
						ReasoningContent: reasoningDelta,
					}:
					case <-ctx.Done():
						return
					}
				}

				prevContent = content
				prevReasoning = reasoning

			case err := <-genErrCh:
				if err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return deltaCh, errCh
}

func Int(v int) *int {
	return &v
}

func Float32(v float32) *float32 {
	return &v
}

func Bool(v bool) *bool {
	return &v
}
