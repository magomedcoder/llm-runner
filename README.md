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

# Запуск (CPU, без CUDA)
make run-cpu

# Запуск (GPU, NVIDIA CUDA)
make run-gpu

# Сборка бинарника (CPU)
make build-cpu

# Сборка бинарника (CUDA)
make build-gpu
```

```bash
./build/llm-runner serve

# Скачать .gguf с Hugging Face
./build/llm-runner download --repo <org/model> --list
./build/llm-runner download --repo <org/model> --file ....gguf

# Клиент к запущенному раннеру
./build/llm-runner remote ping
./build/llm-runner remote models
./build/llm-runner remote run --model mymodel --prompt "привет"
```
