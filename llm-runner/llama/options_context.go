package llama

import (
	"runtime"
)

// WithContext задает размер контекстного окна в токенах
// Чем больше окно, тем больше текста модель учитывает, но выше расход памяти
func WithContext(size int) ContextOption {
	return func(c *contextConfig) {
		c.contextSize = size
	}
}

// WithBatch задает размер батча обработки токенов
// Больший батч может повысить скорость, но также увеличивает нагрузку на память
func WithBatch(size int) ContextOption {
	return func(c *contextConfig) {
		c.batchSize = size
	}
}

// WithThreads задает число CPU-потоков для основных вычислений
// Подбирается под доступные ядра и желаемый баланс скорости/нагрузки
func WithThreads(n int) ContextOption {
	return func(c *contextConfig) {
		c.threads = n
	}
}

// WithThreadsBatch задает число потоков для батч-обработки
// Отдельно настраивается для лучшей производительности на конкретной машине
func WithThreadsBatch(n int) ContextOption {
	return func(c *contextConfig) {
		c.threadsBatch = n
	}
}

// WithF16Memory включает использование f16 для внутренней памяти контекста
// Это снижает потребление памяти ценой возможной потери точности в вычислениях
func WithF16Memory() ContextOption {
	return func(c *contextConfig) {
		c.f16Memory = true
	}
}

// WithEmbeddings включает режим вычисления эмбеддингов
// Используется, когда нужны векторы представления текста вместо генерации ответа
func WithEmbeddings() ContextOption {
	return func(c *contextConfig) {
		c.embeddings = true
	}
}

// WithKVCacheType задает формат хранения KV-кэша (f16, q8_0, q4_0)
// Если значение не из списка, опция игнорируется и текущая конфигурация не меняется
func WithKVCacheType(cacheType string) ContextOption {
	return func(c *contextConfig) {
		switch cacheType {
		case
			"f16",
			"q8_0",
			"q4_0":
			c.kvCacheType = cacheType
		default:
		}
	}
}

// WithFlashAttn задает режим Flash Attention: auto, enabled или disabled
// Неверное значение игнорируется, чтобы не сломать рабочую конфигурацию
func WithFlashAttn(mode string) ContextOption {
	return func(c *contextConfig) {
		switch mode {
		case
			"auto",
			"enabled",
			"disabled":
			c.flashAttn = mode
		default:
		}
	}
}

// WithParallel задает количество параллельных обработчиков запросов
// Минимум принудительно равен 1, чтобы избежать невалидной конфигурации
func WithParallel(n int) ContextOption {
	return func(c *contextConfig) {
		if n < 1 {
			n = 1
		}

		c.nParallel = n
	}
}

// WithPrefixCaching включает или отключает кэширование префикса
// Полезно для ускорения повторяющихся запросов с общим началом промпта
func WithPrefixCaching(enabled bool) ContextOption {
	return func(c *contextConfig) {
		c.prefixCaching = enabled
	}
}

// init выставляет число потоков по умолчанию равным числу cpu-ядер
func init() {
	defaultContextConfig.threads = runtime.NumCPU()
}
