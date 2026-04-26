package llama

type Tool struct {
	Type     string       `json:"type"`     // "function"
	Function ToolFunction `json:"function"` // Описание функции
}

type ToolFunction struct {
	Name        string         `json:"name"`        // Имя функции (должно быть корректным идентификатором)
	Description string         `json:"description"` // Человекочитаемое описание
	Parameters  map[string]any `json:"parameters"`  // json-схема параметров
}

type ToolCall struct {
	ID       string `json:"id"`   // Уникальный идентификатор этого вызова
	Type     string `json:"type"` // "function"
	Function struct {
		Name      string `json:"name"`      // Имя вызываемой функции
		Arguments string `json:"arguments"` // json-строка с аргументами
	} `json:"function"`
}
