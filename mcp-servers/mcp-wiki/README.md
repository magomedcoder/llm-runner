# MCP wiki server

## Transport SSE

Сборка:

```bash
go build -o ./build/mcp-wiki-sse ./mcp-servers/mcp-wiki/cmd/mcp-wiki-sse
```

Запуск:

```bash
./build/mcp-wiki-sse -listen 127.0.0.1:8771 -wiki-dir /путь/к/wiki
```

В GEN:

- `transport = sse`
- `url = http://127.0.0.1:8771/`

Пример: `examples/connection.sse.json`

---

## Transport streamable HTTP

Сборка:

```bash
go build -o ./build/mcp-wiki-streamable ./mcp-servers/mcp-wiki/cmd/mcp-wiki-streamable
```

Запуск:

```bash
./build/mcp-wiki-streamable -listen 127.0.0.1:8772 -wiki-dir /путь/к/wiki
```

В GEN:

- `transport = streamable`
- `url = http://127.0.0.1:8772/`

Пример: `examples/connection.streamable.json`

---

## Инструменты MCP

| Tool                 | Описание                                                                                   |
|----------------------|--------------------------------------------------------------------------------------------|
| `wiki_model_prompts` | Краткий Markdown подсказок для модели (тот же текст, что prompt/resource); без аргументов  |
| `index_wiki_folder`  | Рекурсивная загрузка, извлечение текста, чанки, индекс TF–IDF                              |
| `ask_wiki`           | JSON: цитаты, источники, `reply_style_hint`, фиксированные `note` при отсутствии оснований |
| `ask_wiki_markdown`  | То же в Markdown                                                                           |
| `wiki_index_status`  | Статистика и `retrieval_method`                                                            |
