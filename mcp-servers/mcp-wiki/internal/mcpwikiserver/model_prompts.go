package mcpwikiserver

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const WikiModelPromptsSource = "gen/mcp-servers/mcp-wiki/internal/mcpwikiserver/model_prompts.go (wikiModelInstructionsText)"

const WikiModelPromptsURI = "wiki://mcp-wiki/model_prompts"

const wikiModelInstructionsText = `# MCP Wiki - краткие правила для модели

Сервер **mcp-wiki** (**wiki-rag**). При подключении клиент получает полный список MCP tools и схемы аргументов.

**Инструменты:** wiki_model_prompts (этот текст), index_wiki_folder, wiki_index_status, ask_wiki, ask_wiki_markdown.

**Индексация:** при запуске бинарника с **-wiki-dir** выполняется полная рекурсивная индексация корня и всех вложенных подпапок. После изменений файлов на диске снова вызови **index_wiki_folder** (или incremental: true).

**wiki_index_status** - files, chunks, source_dir. Если chunks == 0 - индекс пуст; сначала **index_wiki_folder**.

**Ответы пользователю** только из **ask_wiki** / **ask_wiki_markdown**: поля answer, sources, note. Не додумывай цитаты и пути к файлам. Если в проиндексированных документах нет оснований - скажи прямо; не заполняй пробелы «общими знаниями».

**Поиск:** непустой **query** (или синонимы из схемы). Не отвечай «нашёл», пока инструмент не вернул данные.

**Ресурс отчёта индекса:** wiki://mcp-wiki/last_index_report
`

func wikiModelPromptsBody() string {
	return strings.TrimSpace(wikiModelInstructionsText)
}

func wikiModelPromptsResourceURIMatch(requestedURI string) bool {
	u := strings.TrimSpace(requestedURI)
	if u == "" {
		return false
	}
	want := WikiModelPromptsURI
	if u == want {
		return true
	}
	return strings.TrimSuffix(u, "/") == strings.TrimSuffix(want, "/")
}

func registerWikiModelPrompts(srv *mcp.Server) {
	if srv == nil {
		return
	}

	body := wikiModelPromptsBody()
	if body == "" {
		return
	}

	srv.AddResource(&mcp.Resource{
		URI:         WikiModelPromptsURI,
		Name:        "model_prompts",
		Description: "Краткий текст подсказок для модели (тот же, что prompt wiki_prompts_full и tool wiki_model_prompts). Правка: " + WikiModelPromptsSource,
		MIMEType:    "text/markdown",
	}, func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req == nil || req.Params == nil {
			return nil, mcp.ResourceNotFoundError("")
		}
		if !wikiModelPromptsResourceURIMatch(req.Params.URI) {
			return nil, mcp.ResourceNotFoundError(strings.TrimSpace(req.Params.URI))
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:      WikiModelPromptsURI,
				MIMEType: "text/markdown",
				Text:     body,
			}},
		}, nil
	})

	srv.AddPrompt(&mcp.Prompt{
		Name:        "wiki_prompts_full",
		Title:       "Подсказки для модели (кратко)",
		Description: "Краткий Markdown: tools, автоиндексация с -wiki-dir, ask_wiki, строгость ответов. Дубликаты: wiki_model_prompts, ресурс " + WikiModelPromptsURI + ". Правка: " + WikiModelPromptsSource + ".",
	}, func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "краткий текст подсказок",
			Messages: []*mcp.PromptMessage{
				{
					Role: "user",
					Content: &mcp.TextContent{
						Text: body,
					},
				},
			},
		}, nil
	})
}
