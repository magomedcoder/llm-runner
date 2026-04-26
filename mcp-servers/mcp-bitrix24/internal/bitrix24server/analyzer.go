package bitrix24server

import (
	"fmt"
	"strings"
	"time"
)

func analyzeTask(task map[string]any, comments []map[string]any, now time.Time) string {
	title := stringField(task, "title", "TITLE")
	status := statusLabel(numberLike(task["status"]))

	createdAt := parseBitrixTime(stringField(task, "createdDate", "CREATED_DATE"))
	deadline := parseBitrixTime(stringField(task, "deadline", "DEADLINE"))

	var out []string
	if title == "" {
		title = "(без названия)"
	}

	out = append(out, fmt.Sprintf("Задача: %s", title))
	out = append(out, fmt.Sprintf("Статус: %s", status))

	if !createdAt.IsZero() {
		out = append(out, fmt.Sprintf("Создана: %s (%d дн. назад)", createdAt.Format(time.RFC3339), int(now.Sub(createdAt).Hours()/24)))
	}

	if deadline.IsZero() {
		out = append(out, "Дедлайн: не указан")
	} else {
		deltaHours := deadline.Sub(now).Hours()
		if deltaHours < 0 {
			out = append(out, fmt.Sprintf("Дедлайн: %s (ПРОСРОЧЕНО на %d ч.)", deadline.Format(time.RFC3339), int(-deltaHours)))
		} else {
			out = append(out, fmt.Sprintf("Дедлайн: %s (осталось ~%d ч.)", deadline.Format(time.RFC3339), int(deltaHours)))
		}
	}

	out = append(out, fmt.Sprintf("Комментариев: %d", len(comments)))
	lastCommentAt := lastCommentTime(comments)
	if !lastCommentAt.IsZero() {
		out = append(out, fmt.Sprintf("Последний комментарий: %s", lastCommentAt.Format(time.RFC3339)))
	}

	risk := "низкий"
	reasons := make([]string, 0, 2)
	if !deadline.IsZero() && deadline.Before(now) {
		risk = "высокий"
		reasons = append(reasons, "задача уже просрочена")
	}
	if len(comments) == 0 && !createdAt.IsZero() && now.Sub(createdAt) > 72*time.Hour {
		if risk == "низкий" {
			risk = "средний"
		}
		reasons = append(reasons, "долгое время нет комментариев")
	}

	if len(reasons) > 0 {
		out = append(out, "Риски: "+risk+" ("+strings.Join(reasons, "; ")+")")
	} else {
		out = append(out, "Риски: "+risk)
	}

	out = append(out, "Рекомендация: проверьте актуальность статуса, ответственного и ближайших шагов.")
	return strings.Join(out, "\n")
}

func statusLabel(status int) string {
	switch status {
	case 1:
		return "Новая"
	case 2:
		return "Ждет выполнения"
	case 3:
		return "Выполняется"
	case 4:
		return "Ожидает контроля"
	case 5:
		return "Завершена"
	case 6:
		return "Отложена"
	case 7:
		return "Отклонена"
	default:
		if status == 0 {
			return "Неизвестно"
		}
		return fmt.Sprintf("Код %d", status)
	}
}

func lastCommentTime(comments []map[string]any) time.Time {
	var latest time.Time
	for _, c := range comments {
		candidate := parseBitrixTime(stringField(c, "post_date", "POST_DATE", "createdDate", "CREATED_DATE", "dateCreate", "DATE_CREATE"))

		if candidate.After(latest) {
			latest = candidate
		}
	}

	return latest
}

func parseBitrixTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}

func stringField(m map[string]any, keys ...string) string {
	for _, k := range keys {
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}

		s, ok := v.(string)
		if ok {
			return strings.TrimSpace(s)
		}
	}

	return ""
}

func numberLike(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case string:
		parsed := strings.TrimSpace(n)
		if parsed == "" {
			return 0
		}

		var value int
		_, err := fmt.Sscanf(parsed, "%d", &value)
		if err == nil {
			return value
		}
	}
	return 0
}
