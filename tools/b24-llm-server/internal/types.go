package internal

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type GenerationWire struct {
	Temperature *float32 `json:"temperature"`
	MaxTokens   *int32   `json:"max_tokens"`
}

type AnalyzeRequest struct {
	TaskID           string         `json:"task_id"`
	TaskTitle        string         `json:"task_title"`
	TaskDescription  string         `json:"task_description"`
	TaskStatus       string         `json:"task_status"`
	TaskDeadline     string         `json:"task_deadline"`
	TaskAssignee     string         `json:"task_assignee"`
	TaskPriority     string         `json:"task_priority"`
	TaskGroupID      string         `json:"task_group_id"`
	TaskCreatedBy    string         `json:"task_created_by"`
	TaskAccomplices  string         `json:"task_accomplices"`
	TaskAuditors     string         `json:"task_auditors"`
	Comments         []TaskComment  `json:"comments"`
	History          []ChatMessage  `json:"history"`
	Prompt           string         `json:"prompt"`
	AnalysisMode     string         `json:"analysis_mode"`
	Generation       *GenerationWire `json:"generation"`
}

type TaskComment struct {
	Author string `json:"author"`
	Text   string `json:"text"`
	Time   string `json:"time"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnalyzeResponse struct {
	Message string `json:"message"`
}

const (
	maxDescriptionRunes = 12000
	maxCommentRunes     = 4000
)

func truncateRunes(s string, max int) string {
	if max <= 0 || s == "" {
		return s
	}
	
	if utf8.RuneCountInString(s) <= max {
		return s
	}

	r := []rune(s)
	return string(r[:max]) + "\n… (текст усечён)"
}

func buildTaskContext(req AnalyzeRequest) string {
	var b strings.Builder
	b.WriteString("Контекст задачи из Bitrix24:\n")
	b.WriteString(fmt.Sprintf("- ID: %s\n", strings.TrimSpace(req.TaskID)))
	b.WriteString(fmt.Sprintf("- Название: %s\n", strings.TrimSpace(req.TaskTitle)))
	b.WriteString(fmt.Sprintf("- Статус: %s\n", strings.TrimSpace(req.TaskStatus)))
	b.WriteString(fmt.Sprintf("- Срок: %s\n", strings.TrimSpace(req.TaskDeadline)))
	b.WriteString(fmt.Sprintf("- Приоритет: %s\n", strings.TrimSpace(req.TaskPriority)))
	b.WriteString(fmt.Sprintf("- Группа/проект (ID): %s\n", strings.TrimSpace(req.TaskGroupID)))
	b.WriteString(fmt.Sprintf("- Исполнитель (ID): %s\n", strings.TrimSpace(req.TaskAssignee)))
	b.WriteString(fmt.Sprintf("- Постановщик (ID): %s\n", strings.TrimSpace(req.TaskCreatedBy)))
	b.WriteString(fmt.Sprintf("- Соисполнители (ID): %s\n", strings.TrimSpace(req.TaskAccomplices)))
	b.WriteString(fmt.Sprintf("- Наблюдатели (ID): %s\n", strings.TrimSpace(req.TaskAuditors)))
	desc := truncateRunes(strings.TrimSpace(req.TaskDescription), maxDescriptionRunes)
	b.WriteString(fmt.Sprintf("- Описание:\n%s\n", desc))
	b.WriteString("\nКомментарии:\n")

	if len(req.Comments) == 0 {
		b.WriteString("- Комментариев нет.\n")
	} else {
		for _, c := range req.Comments {
			text := truncateRunes(strings.TrimSpace(c.Text), maxCommentRunes)
			b.WriteString(fmt.Sprintf("- [%s] %s: %s\n",
				strings.TrimSpace(c.Time),
				strings.TrimSpace(c.Author),
				text,
			))
		}
	}

	return b.String()
}
