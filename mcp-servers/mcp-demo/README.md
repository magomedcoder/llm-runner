### MCP demo - транспорт stdio

Сборка:

```bash
go build -o mcp-stdio-demo ./cmd/mcp-stdio-demo
```

В GEN: transport = stdio, command = абсолютный путь к бинарнику, args пустые

---

### MCP demo - транспорт SSE

Сборка:

```bash
go build -o mcp-sse-demo ./cmd/mcp-sse-demo
```

Запуск:

```bash
./mcp-sse-demo -listen 127.0.0.1:8765
```

В GEN: transport = sse, url = http://127.0.0.1:8765/

Для удалённого хоста добавьте его в mcp.http_allow_hosts (или http_allow_any)

--- 

### MCP demo - транспорт streamable HTTP

Сборка:

```bash
go build -o mcp-streamable-demo ./cmd/mcp-streamable-demo
```

Запуск:

```bash
./mcp-streamable-demo -listen 127.0.0.1:8766
```

В GEN: transport = streamable, url = http://127.0.0.1:8766/
