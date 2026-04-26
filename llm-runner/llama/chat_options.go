package llama

type ChatOptions struct {
	MaxTokens          *int     // Максимальное число токенов для генерации (nil = значение модели по умолчанию)
	Temperature        *float32 // Температура сэмплирования (nil = значение модели по умолчанию, обычно 0.8)
	TopP               *float32 // Порог nucleus-сэмплирования (nil = значение модели по умолчанию, обычно 0.95)
	TopK               *int     // Сэмплирование Top-K (nil = значение модели по умолчанию, обычно 40)
	Seed               *int     // Сид для воспроизводимой генерации (nil = случайный)
	StopWords          []string // Дополнительные стоп-последовательности сверх значений модели по умолчанию
	ChatTemplate       string
	ChatTemplateKwargs map[string]any
	EnableThinking     *bool           // Включить/выключить вывод reasoning (nil = значение модели по умолчанию)
	ReasoningBudget    *int            // Лимит токенов для reasoning (-1 = без ограничений, 0 = отключено)
	ReasoningFormat    ReasoningFormat // Как обрабатывать reasoning-контент
	StreamBufferSize   int             // Размер буфера для потоковых каналов (по умолчанию: 256)
}

type ReasoningFormat int

const (
	ReasoningFormatNone ReasoningFormat = iota
	ReasoningFormatAuto
	ReasoningFormatDeepSeekLegacy
	ReasoningFormatDeepSeek
)

func (r ReasoningFormat) String() string {
	switch r {
	case ReasoningFormatNone:
		return "none"
	case ReasoningFormatAuto:
		return "auto"
	case ReasoningFormatDeepSeekLegacy:
		return "deepseek-legacy"
	case ReasoningFormatDeepSeek:
		return "deepseek"
	default:
		return "unknown"
	}
}
