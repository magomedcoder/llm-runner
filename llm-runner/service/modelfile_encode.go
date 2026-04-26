package service

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseModelYAMLData(data []byte) (*ModelYAML, error) {
	return parseModelYAML(data)
}

func EncodeModelfile(cfg *ModelYAML, weightsBasename string) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("пустой манифест")
	}

	from := strings.TrimSpace(cfg.From)
	if from == "" {
		from = strings.TrimSpace(weightsBasename)
	}
	if from == "" {
		return "", fmt.Errorf("не удалось определить FROM: укажите from в YAML или используйте имя модели с известным файлом весов")
	}

	var b strings.Builder
	fmt.Fprintf(&b, "FROM %s\n", from)

	if cfg.NumCtx != nil && *cfg.NumCtx > 0 {
		fmt.Fprintf(&b, "PARAMETER num_ctx %d\n", *cfg.NumCtx)
	}

	if p := cfg.Parameter; p != nil {
		if p.Temperature != nil {
			fmt.Fprintf(&b, "PARAMETER temperature %s\n", formatFloatModelfile(*p.Temperature))
		}

		if p.TopP != nil {
			fmt.Fprintf(&b, "PARAMETER top_p %s\n", formatFloatModelfile(*p.TopP))
		}

		if p.TopK != nil {
			fmt.Fprintf(&b, "PARAMETER top_k %d\n", *p.TopK)
		}

		if p.MaxTokens != nil {
			fmt.Fprintf(&b, "PARAMETER num_predict %d\n", *p.MaxTokens)
		}

		if p.RepeatLastN != nil {
			fmt.Fprintf(&b, "PARAMETER repeat_last_n %d\n", *p.RepeatLastN)
		}

		if p.RepeatPenalty != nil {
			fmt.Fprintf(&b, "PARAMETER repeat_penalty %s\n", formatFloatModelfile(*p.RepeatPenalty))
		}

		if p.Seed != nil {
			fmt.Fprintf(&b, "PARAMETER seed %d\n", *p.Seed)
		}

		if p.MinP != nil {
			fmt.Fprintf(&b, "PARAMETER min_p %s\n", formatFloatModelfile(*p.MinP))
		}
	}

	for _, s := range cfg.Stop {
		fmt.Fprintf(&b, "PARAMETER stop %s\n", quoteModelfileToken(s))
	}

	if sys := strings.TrimSpace(cfg.System); sys != "" {
		writeModelfileDirectiveBody(&b, "SYSTEM", sys)
	}

	if tmpl := strings.TrimSpace(cfg.Template); tmpl != "" {
		writeModelfileDirectiveBody(&b, "TEMPLATE", tmpl)
	}

	for _, m := range cfg.Messages {
		if strings.TrimSpace(m.Content) == "" {
			continue
		}

		role := strings.TrimSpace(strings.ToLower(m.Role))
		if role == "" {
			role = "user"
		}

		writeModelfileMessage(&b, role, m.Content)
	}

	return b.String(), nil
}

func formatFloatModelfile(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func quoteModelfileToken(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return `""`
	}

	if strings.ContainsAny(s, " \t\n\"") || strings.Contains(s, "#") {
		return strconv.Quote(s)
	}

	return s
}

func writeModelfileDirectiveBody(b *strings.Builder, directive, body string) {
	if strings.Contains(body, "\n") || strings.Contains(body, `"`) {
		fmt.Fprintf(b, "%s \"\"\"\n%s\n\"\"\"\n", directive, body)
		return
	}

	fmt.Fprintf(b, "%s %s\n", directive, body)
}

func writeModelfileMessage(b *strings.Builder, role, content string) {
	if strings.Contains(content, "\n") || strings.Contains(content, `"`) {
		fmt.Fprintf(b, "MESSAGE %s \"\"\"\n%s\n\"\"\"\n", role, content)
		return
	}

	fmt.Fprintf(b, "MESSAGE %s %s\n", role, content)
}
