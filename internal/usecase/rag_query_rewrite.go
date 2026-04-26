package usecase

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/pkg/logger"
)

const ragQueryRewriteSystem = `Ты превращаешь вопрос пользователя в краткий поисковый запрос для семантического поиска по фрагментам документов.

Правила:
- Сохраняй язык сообщения пользователя (в том числе при многоязычном вводе).
- Выведи только текст поискового запроса, одна строка (без кавычек, подписей, ограждений кода).
- Сохраняй ключевые сущности, числа и имена собственные.
- Если ввод уже оптимален, повтори его с минимальными правками.`

const ragQueryRewriteMaxRunes = 512
const ragHyDEMaxRunes = 1200

const ragHyDESystem = `Ты пишешь компактный гипотетический отрывок, в котором с высокой вероятностью содержится ответ на вопрос пользователя.

Правила:
- Тот же язык, что и у пользователя.
- Только простой текст (без markdown, маркированных списков, вступлений).
- Включи важные сущности, термины и ограничения из вопроса.
- 3–8 предложений, фактический стиль, без оговорок.`

func sanitizeRewrittenQuery(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)

	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		first := strings.TrimSpace(s[:idx])
		if first != "" {
			s = first
		}
	}

	return capRAGQueryRunes(s, ragQueryRewriteMaxRunes)
}

func sanitizeHyDEPseudoDocument(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	return capRAGQueryRunes(s, ragHyDEMaxRunes)
}

func capRAGQueryRunes(s string, max int) string {
	if max <= 0 {
		return s
	}

	if utf8.RuneCountInString(s) <= max {
		return s
	}

	r := []rune(s)
	return string(r[:max])
}

func (u *DocumentIngestUseCase) rewriteQueryForRAGRetrieval(ctx context.Context, sessionID int64, userQuery string, chatModel string) (string, error) {
	msgs := []*domain.Message{
		domain.NewMessage(sessionID, ragQueryRewriteSystem, domain.MessageRoleSystem),
		domain.NewMessage(sessionID, userQuery, domain.MessageRoleUser),
	}

	temp := float32(0.15)
	mt := u.queryRewriteMaxTokens
	gp := &domain.GenerationParams{
		Temperature: &temp,
		MaxTokens:   &mt,
	}

	ch, err := u.llmRepo.SendMessage(ctx, sessionID, chatModel, msgs, nil, u.queryRewriteTimeoutSeconds, gp)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case chunk, ok := <-ch:
			if !ok {
				return sanitizeRewrittenQuery(b.String()), nil
			}
			b.WriteString(chunk.Content)
		}
	}
}

func (u *DocumentIngestUseCase) generateHyDEPseudoDocument(ctx context.Context, sessionID int64, query string, chatModel string) (string, error) {
	msgs := []*domain.Message{
		domain.NewMessage(sessionID, ragHyDESystem, domain.MessageRoleSystem),
		domain.NewMessage(sessionID, query, domain.MessageRoleUser),
	}

	temp := float32(0.2)
	mt := u.hydeMaxTokens
	gp := &domain.GenerationParams{
		Temperature: &temp,
		MaxTokens:   &mt,
	}

	ch, err := u.llmRepo.SendMessage(ctx, sessionID, chatModel, msgs, nil, u.hydeTimeoutSeconds, gp)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case chunk, ok := <-ch:
			if !ok {
				out := sanitizeHyDEPseudoDocument(b.String())
				if out != "" {
					logger.I("DocumentIngest: RAG HyDE: псевдодокумент, символов=%d", utf8.RuneCountInString(out))
				}

				return out, nil
			}

			b.WriteString(chunk.Content)
		}
	}
}
