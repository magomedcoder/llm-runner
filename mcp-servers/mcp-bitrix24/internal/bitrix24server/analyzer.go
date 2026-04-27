package bitrix24server

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func analyzeTask(task map[string]any, comments []map[string]any, now time.Time) string {
	title := stringField(task, "title", "TITLE")
	status := statusLabel(numberLike(task["status"]))
	statusCode := numberLike(field(task, "status", "STATUS"))

	createdAt := parseBitrixTime(stringField(task, "createdDate", "CREATED_DATE"))
	changedAt := parseBitrixTime(stringField(task, "changedDate", "CHANGED_DATE"))
	closedAt := parseBitrixTime(stringField(task, "closedDate", "CLOSED_DATE"))
	deadline := parseBitrixTime(stringField(task, "deadline", "DEADLINE"))
	lastActivity := parseBitrixTime(stringField(task, "activityDate", "ACTIVITY_DATE"))

	var out []string
	if title == "" {
		title = "(без названия)"
	}
	taskID := stringField(task, "id", "ID")
	responsible := stringField(task, "responsibleId", "RESPONSIBLE_ID")
	creator := stringField(task, "createdBy", "CREATED_BY")
	priority := numberLike(field(task, "priority", "PRIORITY"))
	timeEstimate := numberLike(field(task, "timeEstimate", "TIME_ESTIMATE"))
	timeSpent := numberLike(field(task, "timeSpentInLogs", "TIME_SPENT_IN_LOGS"))

	out = append(out, "=== Паспорт задачи ===")
	if taskID != "" {
		out = append(out, fmt.Sprintf("ID: %s", taskID))
	}

	out = append(out, fmt.Sprintf("Название: %s", title))
	out = append(out, fmt.Sprintf("Статус: %s (код %d)", status, statusCode))
	if creator != "" || responsible != "" {
		out = append(out, fmt.Sprintf("Постановщик: %s | Ответственный: %s", emptyDash(creator), emptyDash(responsible)))
	}

	if priority > 0 {
		out = append(out, fmt.Sprintf("Приоритет: %d", priority))
	}

	if !createdAt.IsZero() {
		out = append(out, fmt.Sprintf("Создана: %s (%d дн. назад)", createdAt.Format(time.RFC3339), int(now.Sub(createdAt).Hours()/24)))
	}

	if !changedAt.IsZero() {
		out = append(out, fmt.Sprintf("Последнее изменение: %s", changedAt.Format(time.RFC3339)))
	}

	if !lastActivity.IsZero() {
		out = append(out, fmt.Sprintf("Последняя активность: %s", lastActivity.Format(time.RFC3339)))
	}

	if !closedAt.IsZero() {
		out = append(out, fmt.Sprintf("Закрыта: %s", closedAt.Format(time.RFC3339)))
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

	if timeEstimate > 0 || timeSpent > 0 {
		execLine := fmt.Sprintf("Трудозатраты: потрачено %d сек.", timeSpent)
		if timeEstimate > 0 {
			execLine += fmt.Sprintf(" / оценка %d сек.", timeEstimate)
		}

		if timeEstimate > 0 {
			execLine += fmt.Sprintf(" (%d%%)", int(float64(timeSpent)*100.0/float64(max(1, timeEstimate))))
		}
		out = append(out, execLine)
	}

	out = append(out, "")
	out = append(out, "=== Анализ коммуникации ===")
	out = append(out, fmt.Sprintf("Комментариев: %d", len(comments)))
	lastCommentAt := lastCommentTime(comments)
	if !lastCommentAt.IsZero() {
		out = append(out, fmt.Sprintf("Последний комментарий: %s (%d ч. назад)", lastCommentAt.Format(time.RFC3339), int(now.Sub(lastCommentAt).Hours())))
	}

	authors := commentAuthors(comments)
	if len(authors) > 0 {
		out = append(out, fmt.Sprintf("Участников в комментариях: %d (%s)", len(authors), strings.Join(authors, ", ")))
	}

	if len(comments) > 0 {
		out = append(out, "Последние комментарии:")
		for _, line := range recentCommentSnippets(comments, 3) {
			out = append(out, "- "+line)
		}
	}

	risk := "низкий"
	score := 0
	reasons := make([]string, 0, 6)
	if !deadline.IsZero() && deadline.Before(now) {
		score += 4
		reasons = append(reasons, "задача уже просрочена")
	}

	if deadline.IsZero() && (statusCode == 2 || statusCode == 3 || statusCode == 4) {
		score += 2
		reasons = append(reasons, "нет дедлайна у активной задачи")
	}

	if len(comments) == 0 && !createdAt.IsZero() && now.Sub(createdAt) > 48*time.Hour {
		score += 2
		reasons = append(reasons, "долгое время нет комментариев")
	}

	if !lastCommentAt.IsZero() && now.Sub(lastCommentAt) > 72*time.Hour && statusCode != 5 && statusCode != 7 {
		score += 2
		reasons = append(reasons, "нет свежих коммуникаций по задаче")
	}

	if scoreMentionsBlockers(comments) {
		score += 2
		reasons = append(reasons, "в комментариях есть сигналы блокировки/ожидания")
	}

	if timeEstimate > 0 && timeSpent > timeEstimate {
		score += 2
		reasons = append(reasons, "превышена оценка времени")
	}

	if statusCode == 5 || statusCode == 7 {
		score = 0
		reasons = nil
	}

	switch {
	case score >= 5:
		risk = "высокий"
	case score >= 2:
		risk = "средний"
	default:
		risk = "низкий"
	}

	out = append(out, "")
	out = append(out, "=== Риски и рекомендации ===")
	if len(reasons) > 0 {
		out = append(out, "Риск: "+risk+" ("+strings.Join(reasons, "; ")+")")
	} else {
		out = append(out, "Риск: "+risk)
	}

	for _, action := range recommendActions(task, comments, now, deadline, lastCommentAt, risk) {
		out = append(out, "- "+action)
	}

	if len(reasons) == 0 && (statusCode == 5 || statusCode == 7) {
		out = append(out, "- Задача уже в финальном статусе, проверьте только корректность итогов и закрывающих комментариев.")
	}

	out = append(out, "")
	out = append(out, "=== Вывод по задаче ===")
	out = append(out, taskConclusion(title, statusCode, risk, deadline, now, len(comments), scoreMentionsBlockers(comments)))

	return strings.Join(out, "\n")
}

func field(m map[string]any, keys ...string) any {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			return v
		}
	}

	return nil
}

func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}

	return v
}

func commentAuthors(comments []map[string]any) []string {
	seen := map[string]struct{}{}
	for _, c := range comments {
		id := stringField(c, "authorId", "AUTHOR_ID")
		if id == "" {
			continue
		}
		seen[id] = struct{}{}
	}

	if len(seen) == 0 {
		return nil
	}

	out := make([]string, 0, len(seen))
	for id := range seen {
		out = append(out, id)
	}
	sort.Strings(out)

	return out
}

func recentCommentSnippets(comments []map[string]any, limit int) []string {
	if limit <= 0 || len(comments) == 0 {
		return nil
	}

	type item struct {
		at   time.Time
		text string
	}

	buf := make([]item, 0, len(comments))
	for _, c := range comments {
		at := parseBitrixTime(stringField(c, "post_date", "POST_DATE", "createdDate", "CREATED_DATE", "dateCreate", "DATE_CREATE"))
		text := strings.TrimSpace(stringField(c, "postMessage", "POST_MESSAGE", "message", "MESSAGE"))
		if text == "" {
			continue
		}

		text = strings.ReplaceAll(text, "\n", " ")
		if len([]rune(text)) > 90 {
			r := []rune(text)
			text = string(r[:90]) + "..."
		}

		buf = append(buf, item{at: at, text: text})
	}

	sort.SliceStable(buf, func(i, j int) bool { return buf[i].at.After(buf[j].at) })

	if len(buf) > limit {
		buf = buf[:limit]
	}

	out := make([]string, 0, len(buf))
	for _, v := range buf {
		if v.at.IsZero() {
			out = append(out, v.text)
			continue
		}

		out = append(out, fmt.Sprintf("%s: %s", v.at.Format(time.RFC3339), v.text))
	}

	return out
}

func scoreMentionsBlockers(comments []map[string]any) bool {
	keywords := []string{
		"блок", "жду", "ожида", "проблем", "не могу", "завис", "risk", "blocked", "blocker",
	}

	for _, c := range comments {
		text := strings.ToLower(strings.TrimSpace(stringField(c, "postMessage", "POST_MESSAGE", "message", "MESSAGE")))
		if text == "" {
			continue
		}

		for _, kw := range keywords {
			if strings.Contains(text, kw) {
				return true
			}
		}
	}

	return false
}

func recommendActions(task map[string]any, comments []map[string]any, now, deadline, lastCommentAt time.Time, risk string) []string {
	statusCode := numberLike(field(task, "status", "STATUS"))
	var actions []string
	if !deadline.IsZero() && deadline.Before(now) && statusCode != 5 && statusCode != 7 {
		actions = append(actions, "Пересогласуйте срок или срочно обновите план закрытия с ответственным.")
	}

	if deadline.IsZero() && statusCode != 5 && statusCode != 7 {
		actions = append(actions, "Добавьте дедлайн, чтобы задача попала в контролируемый контур.")
	}

	if len(comments) == 0 {
		actions = append(actions, "Запросите статус-апдейт у ответственного одним комментарием с датой следующего шага.")
	}

	if !lastCommentAt.IsZero() && now.Sub(lastCommentAt) > 72*time.Hour && statusCode != 5 && statusCode != 7 {
		actions = append(actions, "Обновите коммуникацию: уточните блокеры и зафиксируйте следующий контрольный чекпоинт.")
	}

	if scoreMentionsBlockers(comments) {
		actions = append(actions, "Разберите блокеры из обсуждения: назначьте владельца каждой проблемы и дедлайн на снятие.")
	}

	if risk == "низкий" {
		actions = append(actions, "Оставьте текущий ритм: краткий weekly-апдейт и контроль дедлайна.")
	}

	if len(actions) == 0 {
		actions = append(actions, "Проверьте актуальность статуса, ответственного и ближайших шагов.")
	}

	return actions
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func taskConclusion(title string, statusCode int, risk string, deadline, now time.Time, commentsCount int, hasBlockers bool) string {
	title = emptyDash(strings.TrimSpace(title))
	if statusCode == 5 || statusCode == 7 {
		return fmt.Sprintf("Задача \"%s\" находится в финальном статусе; риск %s. Действие: проверить итог и закрывающие артефакты.", title, risk)
	}

	if !deadline.IsZero() && deadline.Before(now) {
		if hasBlockers {
			return fmt.Sprintf("Задача \"%s\" просрочена и имеет признаки блокировки; риск %s. Действие: срочная эскалация, снять блокер и зафиксировать новый план.", title, risk)
		}
		return fmt.Sprintf("Задача \"%s\" просрочена; риск %s. Действие: срочно обновить план работ и согласовать новый реалистичный дедлайн.", title, risk)
	}

	if commentsCount == 0 {
		return fmt.Sprintf("Задача \"%s\" без коммуникаций в комментариях; риск %s. Действие: запросить статус-апдейт и ближайший следующий шаг.", title, risk)
	}

	if hasBlockers {
		return fmt.Sprintf("По задаче \"%s\" есть сигналы блокера; риск %s. Действие: назначить владельца блокера и срок снятия.", title, risk)
	}

	if deadline.IsZero() {
		return fmt.Sprintf("Задача \"%s\" активна, но без дедлайна; риск %s. Действие: установить контрольный срок и точку проверки.", title, risk)
	}

	return fmt.Sprintf("Задача \"%s\" контролируема; риск %s. Действие: поддерживать текущий ритм обновлений и мониторинг дедлайна.", title, risk)
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
