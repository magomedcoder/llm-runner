# Tce-server

Лёгкий HTTP-адаптер для интеграции Tce с gen-раннером, с api в формате JSON/SSE

Формат: JSON, для стриминга - SSE (`text/event-stream`)

## Запуск

```bash
go run ./tools/tce-server/cmd
```

Переменные окружения:

- `PORT` - HTTP порт сервера
- `GEN_RUNNER_ADDR` - адрес gen-раннера
- `GEN_MODEL` - идентификатор модели для запросов

## 1) Проверка доступности

### `GET /tce/v1/health`

Назначение: проверка доступности сервиса и связности с gen-раннером

Успешный ответ `200`:

```json
{
  "ok": true
}
```

Если раннер недоступен: `503` + `error.code = "internal_error"`

## 2) Чат

### `POST /tce/v1/chat`

Тело запроса:

```json
{
  "stream": true,
  "system": "Ты помощник по коду.",
  "messages": [
    {
      "role": "user",
      "content": "Объясни этот фрагмент"
    }
  ],
  "editor": {
    "path": "src/main.rs",
    "language": "rust",
    "snippet": "fn main() {}",
    "cursor_line": 0,
    "cursor_column": 0
  },
  "generate": {
    "max_tokens": 1024,
    "temperature": 0.2
  }
}
```

Поля запроса:

- `stream` (`boolean`, обязательно) - `true` для SSE, `false` для единого JSON-ответа
- `system` (`string`, обязательно) - системная инструкция; может быть пустой строкой
- `messages` (`array`, обязательно) - история диалога; последний элемент должен иметь `role: "user"`
- `editor` (`object`, опционально) - контекст активного файла
- `generate` (`object`, опционально) - параметры генерации (`max_tokens`, `temperature`)

### Ответ без стрима (`stream: false`)

Ответ `200`:

```json
{
  "message": {
    "role": "assistant",
    "content": "Ответ модели"
  },
  "finish": "stop"
}
```

### Ответ со стримом (`stream: true`, SSE)

Сервис отправляет события:

```text
event: delta
data: {"text":"часть ответа"}

event: delta
data: {"text":"еще часть"}

event: end
data: {"finish":"stop"}
```

Правила:

- Событие `delta` содержит JSON с обязательным полем `text`
- Финальное событие `end` содержит поле `finish`
- Конкатенация всех значений `text` из `delta` формирует полный ответ ассистента

## 3) Шаг агента

### `POST /tce/v1/agent/step`

Назначение: выполнение одного шага автономного агента

Эндпоинт выполняет **один шаг** агента: принимает состояние сессии и observations от предыдущих инструментов, возвращает список tool-вызовов или завершение.

#### Запрос

Минимальный совместимый формат:

```json
{
  "session_id": "session-1",
  "goal": "Починить тесты в модуле parser",
  "context": {
    "workspace_root": ".",
    "branch": "feature/agent"
  },
  "observations": [
    {
      "call_id": "call-1",
      "tool": "read_file",
      "ok": true,
      "result": {
        "path": "src/main.rs",
        "content": "fn main() {}"
      }
    }
  ]
}
```

Поля запроса:

- `session_id` (`string`, рекомендуется) - идентификатор сессии агента
- `goal` (`string`, рекомендуется) - цель пользователя для текущей сессии
- `context` (`object`, опционально) - мета-контекст шага (workspace, branch, режим)
- `observations` (`array`, опционально) - результаты исполнения инструментов на предыдущем шаге

Текущая реализация сервера допускает запрос без полей (stub-режим), но для рабочего цикла рекомендуется формат выше.

#### Успешный ответ

Минимальный формат успешного ответа:

```json
{
  "finish": false,
  "summary": "Проверил файлы и подготовил следующие действия.",
  "calls": [
    {
      "tool": "read_file",
      "id": "call-1",
      "args": {
        "path": "src/main.rs"
      }
    }
  ]
}
```

Поля ответа:

- `finish` (`boolean`, обязательно) - завершить ли сессию
- `summary` (`string`, обязательно) - краткое описание результата шага для UI
- `calls` (`array`, обязательно, может быть пустым) - список tool-вызовов на исполнение

Если `finish=true`, клиент может завершить цикл агента и показать итог `summary`.

#### Формат вызова инструмента

```json
{
  "tool": "apply_patch",
  "id": "call-2",
  "args": {
    "path": "src/lib.rs",
    "patch": "*** Begin Patch ..."
  }
}
```

Поля:

- `tool` (`string`, обязательно) - имя инструмента
- `id` (`string`, обязательно) - уникальный id вызова в рамках шага
- `args` (`object`, обязательно) - аргументы инструмента

#### Рекомендуемые инструменты (MVP)

- `list_dir`
- `read_file`
- `glob_search`
- `search_content`
- `apply_patch`
- `create_file`
- `run_command`
- `finish` (как отдельный `tool` опционально, либо через `finish=true`)

#### Совместимость с текущей реализацией

Текущий `tce-server` возвращает stub-ответ с `finish/summary/calls`; этого достаточно для первичной клиентской интеграции.

## 4) Ошибки

При ошибке сервис возвращает HTTP `4xx/5xx` и тело JSON:

```json
{
  "error": {
    "code": "bad_request",
    "message": "Описание ошибки"
  }
}
```

Используемые `error.code`:

- `bad_request`
- `method_not_allowed`
- `internal_error`
