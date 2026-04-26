# MCP Bitrix24

Сервер MCP для задач Bitrix24 через входящий webhook:

- список задач (`tasks.task.list`)
- получение задачи (`tasks.task.get`)
- комментарии (`tasks.task.commentitem.getlist`)
- сводка по задаче (`b24_analyze_task`)
- произвольный REST-вызов (`b24_call_method`)

## Переменные окружения

| Переменная         | Назначение                                                                       |
|--------------------|----------------------------------------------------------------------------------|
| `B24_WEBHOOK_BASE` | Базовый URL webhook, например `https://bitrix24.example.com/rest/43176/00000000` |

---

## Transport stdio

Сборка:

```bash
go build -o ./build/mcp-bitrix24-stdio ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-stdio
```

В GEN:

- `transport = stdio`
- `command` - путь к бинарнику
- `args` - обычно пустой массив (URL задаётся через `B24_WEBHOOK_BASE`)

---

## Transport SSE

Сборка:

```bash
go build -o ./build/mcp-bitrix24-sse ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-sse
```

Запуск:

```bash
export B24_WEBHOOK_BASE="https://bitrix24.example.com/rest/43176/00000000"
./build/mcp-bitrix24-sse -listen 127.0.0.1:8785
```

```
transport = sse

url = http://127.0.0.1:8785/
```

---

## Transport streamable HTTP

Сборка:

```bash
go build -o ./build/mcp-bitrix24-streamable ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-streamable
```

Запуск:

```bash
export B24_WEBHOOK_BASE="https://bitrix24.example.com/rest/43176/00000000"
./build/mcp-bitrix24-streamable -listen 127.0.0.1:8786
```

```
transport = streamable

url = http://127.0.0.1:8786/
```

---

## Инструменты MCP

Схемы аргументов отдаёт сам сервер MCP (поля инструментов):

| Tool                    | Назначение                                                      |
|-------------------------|-----------------------------------------------------------------|
| `b24_list_tasks`        | `tasks.task.list` (`filter`, `select`, `order`, `start`)        |
| `b24_get_task`          | `tasks.task.get` (`task_id`, `select`)                          |
| `b24_get_task_comments` | `tasks.task.commentitem.getlist` (`task_id`, `order`, `select`) |
| `b24_analyze_task`      | Сводка по задаче и комментариям (`task_id`, `include_comments`) |
| `b24_call_method`       | Произвольный REST-метод (`method`, `params`)                    |

---

## Mock REST (локальная отладка)

Отдельный бинарник - не MCP, а простой HTTP-мок Bitrix REST:

```bash
go build -o ./build/mcp-bitrix24-mock-rest ./mcp-servers/mcp-bitrix24/cmd/mcp-bitrix24-mock-rest
./build/mcp-bitrix24-mock-rest -listen 127.0.0.1:8899
```

База методов: `http://127.0.0.1:8899/rest/43176/mock-token/<method>`.
