# Сборка клиента

### Зависимости

- Flutter 3.24+
- Dart SDK ^3.10.7

## Поддерживаемые платформы

| Платформа | Версия                                |
|-----------|---------------------------------------|
| Linux     | glibc 2.31+ (Ubuntu 20.04+ и аналоги) |
| Android   | 7.0+                                  |
| iOS       | 13.0+                                 |
| macOS     | Catalina 10.15+                       |
| Windows   | 10+                                   |

## Linux и Android (Docker)

Сборка **Linux** и **Android** через Docker:

```bash
docker build -f Dockerfile-client-app --target linux-build -t gen-app-linux .
docker run --rm -e TARGETS=linux,android -v ./out:/opt/gen/out gen-app-linux
```

## Windows (Docker на хосте Windows)

Сборка **Windows** возможна только на **хосте Windows**:

```bash
docker build -f Dockerfile-client-app-windows --target windows-build -t gen-app-windows .
docker run --rm -v .\out:C:\gen\out gen-app-windows
```
