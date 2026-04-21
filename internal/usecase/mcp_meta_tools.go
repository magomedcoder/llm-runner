package usecase

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/pkg/logger"
)

func (c *ChatUseCase) appendMCPLLMContext(ctx context.Context, msg *domain.Message, settings *domain.ChatSessionSettings, userID int) {
	if msg == nil || settings == nil || !settings.MCPEnabled || c.mcpServerRepo == nil {
		return
	}
	effective := c.mcpEffectiveServerIDs(ctx, userID, settings)
	if len(effective) == 0 {
		return
	}

	logger.D("MCP appendMCPLLMContext: user_id=%d разрешённых_server_id=%d", userID, len(effective))
	var b strings.Builder
	b.WriteString("[MCP] В этой сессии чата включены внешние инструменты. Разрешённые server_id (используй только их):\n")
	for _, sid := range effective {
		if sid <= 0 {
			continue
		}
		line := fmt.Sprintf("- id=%d", sid)
		if srv, err := c.mcpServerRepo.GetByIDAccessible(ctx, sid, userID); err == nil && srv != nil {
			if srv.Enabled {
				if n := strings.TrimSpace(srv.Name); n != "" {
					line = fmt.Sprintf("- id=%d · %s", sid, n)
				}
			} else {
				line = fmt.Sprintf("- id=%d · (отключён в каталоге)", sid)
			}
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}

	b.WriteString("\nИмена MCP-инструментов вида mcp_<id>_h<hex> задаёт платформа; hex привязан к реальному имени на сервере. Точные имена для этого запроса - в блоке [Tools] ниже (если он есть) и в payload раннера; они совпадают.\n")
	b.WriteString("КРИТИЧНО - имена tools:\n")
	b.WriteString("- Копируй поле name из блока [Tools] / списка tools СИМВОЛ В СИМВОЛ. Не сокращай, не «улучшай», не подставляй примеры вроде mcp_1_h123456 или шаблонные hex.\n")
	b.WriteString("- Любое другое имя (включая похожее на mcp_<id>_h...) не будет выполнено: платформа не угадает вашу замену.\n")
	b.WriteString("КРИТИЧНО - как выполняется вызов:\n")
	b.WriteString("- Недостаточно описать вызов в свободном тексте («предположу инструмент…», «если вернётся…»). Чтобы инструмент реально вызвался, в ответе должен быть машиночитаемый вызов в формате, ожидаемом для tool-calling (JSON-массив с полями tool_name и parameters и/или блок ```json … ``` - как в вашей инструкции к модели).\n")
	b.WriteString("- Сначала вызови релевантный tool с корректными аргументами по его JSON-схеме, получи данные, затем формируй ответ пользователю по факту результата.\n")
	b.WriteString("Не добавляй в аргументы поле server_id: привязка к серверу уже зашита в имени инструмента.\n")
	b.WriteString("Не утверждай, что инструмента нет или что доступ невозможен, пока не проверишь доступные tools и не выполнил релевантный вызов.\n")

	msg.Content += "\n\n" + strings.TrimSpace(b.String())
}

func (c *ChatUseCase) appendResolvedToolCatalog(msg *domain.Message, genParams *domain.GenerationParams) {
	if msg == nil || genParams == nil || len(genParams.Tools) == 0 {
		return
	}

	var b strings.Builder
	b.WriteString("[Tools] Разрешённые инструменты в этом запросе - в вызовах используй только эти значения name (символ в символ):\n")
	for _, t := range genParams.Tools {
		name := strings.TrimSpace(t.Name)
		if name == "" {
			name = "(без имени)"
		}
		b.WriteString("- ")
		b.WriteString(name)
		b.WriteByte('\n')
	}

	text := strings.TrimSpace(b.String())
	if text == "" {
		return
	}

	if utf8.RuneCountInString(text) > maxLLMToolNamesListRunes {
		runes := []rune(text)
		text = string(runes[:maxLLMToolNamesListRunes]) + fmt.Sprintf("\n…(обрезано, всего инструментов=%d)", len(genParams.Tools))
	}

	msg.Content += "\n\n" + text
}
