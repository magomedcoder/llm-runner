# Сборка сервера

### Зависимости

- Go 1.26+
- Protobuf (`protoc`) 30.2+
- `make`
- PostgreSQL 16+

## Локальная сборка

```bash
go mod download
make install
make build
```

Команда `make build` собирает серверный бинарник в `./build/gen-server`

Сервер использует конфигурацию из `configs/config.yaml` (можно взять за основу `configs/config.example.yaml`).

## Запуск сервера в dev-режиме

```bash
make run
```
