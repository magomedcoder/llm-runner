package internal

import "strings"

func systemPromptForAnalysisMode(mode string) string {
	base := "Ты аналитик задач Bitrix24. Отвечай структурировано и по делу. Используй контекст задачи и комментариев."
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "summarize":
		return base + " Режим: краткое резюме (3-5 предложений) для руководителя."
	case "plan":
		return base + " Режим: план работ - нумерованный список шагов, порядок и зависимости."
	case "risks":
		return base + " Режим: риски и блокеры - список с пояснением серьёзности."
	case "clarify":
		return base + " Режим: уточняющие вопросы постановщику - маркированный список."
	case "acceptance":
		return base + " Режим: критерии приёмки - измеримые пункты."
	case "draft_comment":
		return "Ты помощник по задачам Bitrix24. Напиши один связный комментарий от лица исполнителя: вежливо, по существу, готовый к публикации."
	default:
		return base + " Отмечай риски, шаги и уточняющие вопросы, где уместно."
	}
}

func effectiveUserPrompt(req AnalyzeRequest) string {
	p := strings.TrimSpace(req.Prompt)
	if p != "" {
		return p
	}

	switch strings.ToLower(strings.TrimSpace(req.AnalysisMode)) {
	case "summarize":
		return "Сделай краткое резюме задачи для руководителя (3–5 предложений)."
	case "plan":
		return "Составь пошаговый план работ (WBS) с учётом контекста и комментариев."
	case "risks":
		return "Выдели риски, блокеры и неясности постановки; предложи смягчение рисков."
	case "clarify":
		return "Сформулируй уточняющие вопросы постановщику, чтобы закрыть пробелы в ТЗ."
	case "acceptance":
		return "Предложи чёткие критерии приёмки (acceptance criteria) по текущему описанию и обсуждению."
	case "draft_comment":
		return "Подготовь нейтральный черновик комментария к задаче для публикации в Bitrix24."
	default:
		return ""
	}
}
