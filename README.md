# LLM Runner

#### Сборка

```bash
# Установка необходимых зависимостей и клонирование llama.cpp
make deps

# Генерация proto
make gen

# Сборка libllama.a (без CUDA)
make build-llama

# Сборка libllama.a с поддержкой NVIDIA (CUDA)
make build-llama-cublas

# Запуск
make run

# Сборка бинарника (CUDA)
make build-nvidia
```

### Скачивание моделей с Hugging Face

```bash
make build-download-model

# Список доступных .gguf в репозитории
./build/download-model -repo ... -list

# Скачать один файл
./build/download-model -repo ... -file ...gguf

# Скачать все .gguf из репозитория
./build/download-model -repo ...
```