package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseModelfile(r io.Reader) (*ModelYAML, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var out ModelYAML
	var param ModelYAMLParameter
	haveParam := false

	for i := 0; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		dir, arg := splitModelfileDirective(trimmed)
		if dir == "" {
			return nil, fmt.Errorf("modelfile: строка %d: неизвестная или некорректная директива: %q", i+1, trimmed)
		}

		switch dir {
		case "LICENSE":
			if strings.HasPrefix(arg, `"""`) {
				_, j, err := consumeTripleQuoted(lines, i, arg)
				if err != nil {
					return nil, fmt.Errorf("modelfile: %w", err)
				}
				i = j
			}
			continue

		case "ADAPTER":
			return nil, fmt.Errorf("modelfile: директива ADAPTER не поддерживается текущим API llama")

		case "REQUIRES":
			continue

		case "MESSAGE":
			role, content, j, err := parseMessageDirective(lines, i, arg)
			if err != nil {
				return nil, err
			}

			i = j
			out.Messages = append(out.Messages, ModelYAMLMessage{Role: role, Content: content})

		case "FROM":
			if arg == "" {
				return nil, fmt.Errorf("modelfile: FROM без значения (строка %d)", i+1)
			}

			if out.From != "" {
				return nil, fmt.Errorf("modelfile: повторная директива FROM")
			}

			out.From = unquoteOllamaPath(arg)

		case "SYSTEM":
			text, j, err := parseModelfileStringField(lines, i, arg, "SYSTEM")
			if err != nil {
				return nil, err
			}

			i = j
			if out.System != "" {
				return nil, fmt.Errorf("modelfile: повторная директива SYSTEM")
			}

			out.System = text

		case "TEMPLATE":
			text, j, err := parseModelfileStringField(lines, i, arg, "TEMPLATE")
			if err != nil {
				return nil, err
			}

			i = j
			if out.Template != "" {
				return nil, fmt.Errorf("modelfile: повторная директива TEMPLATE")
			}

			out.Template = text

		case "PARAMETER":
			key, val, err := parseParameterArg(arg)
			if err != nil {
				return nil, fmt.Errorf("modelfile: PARAMETER: %w", err)
			}

			if err := applyParameter(&out, &param, &haveParam, key, val); err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("modelfile: строка %d: неизвестная директива %q", i+1, dir)
		}
	}

	if strings.TrimSpace(out.From) == "" {
		return nil, fmt.Errorf("modelfile: требуется директива FROM")
	}

	if haveParam {
		out.Parameter = &param
	}

	return &out, nil
}

func parseModelfileStringField(lines []string, i int, arg, name string) (string, int, error) {
	if !strings.HasPrefix(arg, `"""`) {
		return arg, i, nil
	}

	body, j, err := consumeTripleQuoted(lines, i, arg)
	if err != nil {
		return "", i, fmt.Errorf("modelfile: %s: %w", name, err)
	}

	return body, j, nil
}

func consumeTripleQuoted(lines []string, start int, firstArg string) (body string, endIndex int, err error) {
	rest := strings.TrimPrefix(firstArg, `"""`)
	if strings.HasSuffix(rest, `"""`) && len(rest) >= 3 {
		return strings.TrimRight(strings.TrimSuffix(rest, `"""`), "\r\n"), start, nil
	}

	var b strings.Builder
	if rest != "" {
		b.WriteString(rest)
	}

	for j := start + 1; j < len(lines); j++ {
		line := strings.TrimRight(lines[j], "\r")
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, `"""`) {
			prefix := line
			if idx := strings.LastIndex(prefix, `"""`); idx >= 0 {
				prefix = prefix[:idx]
			}
			if b.Len() > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(strings.TrimSuffix(prefix, `"""`))

			return strings.TrimRight(b.String(), "\r\n"), j, nil
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}

		b.WriteString(line)
	}

	return "", start, fmt.Errorf("не закрыта тройная кавычка \"\"\"")
}

func parseParameterArg(arg string) (key, val string, err error) {
	fields := strings.Fields(arg)
	if len(fields) < 2 {
		return "", "", fmt.Errorf("ожидалось PARAMETER <имя> <значение>, получено %q", arg)
	}

	key = strings.ToLower(fields[0])
	rest := strings.TrimSpace(arg[len(fields[0]):])
	val = unquoteParameterValue(rest)

	return key, val, nil
}

func unquoteParameterValue(s string) string {
	s = strings.TrimSpace(s)
	if u, err := strconv.Unquote(s); err == nil {
		return u
	}

	return s
}

func applyParameter(out *ModelYAML, param *ModelYAMLParameter, have *bool, key, val string) error {
	switch key {
	case "temperature":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("temperature: %w", err)
		}

		param.Temperature = &f
		*have = true
	case "top_p":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("top_p: %w", err)
		}

		param.TopP = &f
		*have = true
	case "top_k":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("top_k: %w", err)
		}

		param.TopK = &n
		*have = true
	case "max_tokens", "num_predict":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}

		param.MaxTokens = &n
		*have = true
	case "num_ctx":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("num_ctx: %w", err)
		}

		if n <= 0 {
			return fmt.Errorf("num_ctx: ожидалось положительное число")
		}

		out.NumCtx = &n
	case "stop":
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("stop: пустое значение")
		}

		out.Stop = append(out.Stop, val)
	case "repeat_last_n":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("repeat_last_n: %w", err)
		}

		param.RepeatLastN = &n
		*have = true
	case "repeat_penalty":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("repeat_penalty: %w", err)
		}

		param.RepeatPenalty = &f
		*have = true
	case "seed":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("seed: %w", err)
		}

		param.Seed = &n
		*have = true
	case "min_p":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("min_p: %w", err)
		}

		param.MinP = &f
		*have = true
	default:
		return fmt.Errorf("неподдерживаемый PARAMETER %q", key)
	}

	return nil
}

func parseMessageDirective(lines []string, lineIdx int, arg string) (role, content string, endIdx int, err error) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return "", "", lineIdx, fmt.Errorf("MESSAGE: не указана роль")
	}

	fields := strings.Fields(arg)
	role = strings.ToLower(fields[0])
	switch role {
	case "system", "user", "assistant":
	default:
		return "", "", lineIdx, fmt.Errorf("MESSAGE: неизвестная роль %q (допустимы system, user, assistant)", role)
	}

	rest := strings.TrimSpace(arg[len(fields[0]):])
	if strings.HasPrefix(rest, `"""`) {
		body, j, e := parseModelfileStringField(lines, lineIdx, rest, "MESSAGE")
		if e != nil {
			return "", "", lineIdx, e
		}

		return role, body, j, nil
	}

	return role, rest, lineIdx, nil
}

func splitModelfileDirective(trimmed string) (dir, arg string) {
	i := strings.IndexFunc(trimmed, unicodeSpace)
	if i < 0 {
		return strings.ToUpper(trimmed), ""
	}

	return strings.ToUpper(trimmed[:i]), strings.TrimSpace(trimmed[i+1:])
}

func unicodeSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func unquoteOllamaPath(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.TrimSpace(s[1 : len(s)-1])
	}

	return s
}

func PrepareModelYAMLForDisk(modelsDir, stem string, cfg *ModelYAML) (*ModelYAML, error) {
	if cfg == nil {
		return nil, fmt.Errorf("пустой манифест")
	}

	if strings.TrimSpace(stem) == "" {
		return nil, fmt.Errorf("пустое имя модели")
	}

	if strings.ContainsAny(stem, `/\`) {
		return nil, fmt.Errorf("имя модели не должно содержать путь: %q", stem)
	}

	cfg, err := cloneNormalizeManifestRefs(cfg, modelsDir)
	if err != nil {
		return nil, err
	}

	fromTrim := strings.TrimSpace(cfg.From)
	if fromTrim == "" {
		return nil, fmt.Errorf("в манифесте нет FROM")
	}

	resolvedFrom, err := ResolveGGUFFile(modelsDir, fromTrim)
	if err != nil {
		return nil, fmt.Errorf("FROM %q: %w", fromTrim, err)
	}

	stemGGUF := filepath.Join(modelsDir, stem+".gguf")
	if st, e := os.Stat(stemGGUF); e == nil && !st.IsDir() {
		if !strings.EqualFold(resolvedFrom, filepath.Base(stemGGUF)) {
			return nil, fmt.Errorf("в каталоге уже есть %s.gguf, а FROM указывает на другой файл (%q); выберите другое имя манифеста", stem, resolvedFrom)
		}
		cp := *cfg
		cp.From = ""

		return &cp, nil
	}

	cp := *cfg
	cp.From = resolvedFrom

	return &cp, nil
}

func WriteModelManifest(modelsDir, stem string, cfg *ModelYAML, overwrite bool) error {
	prepared, err := PrepareModelYAMLForDisk(modelsDir, stem, cfg)
	if err != nil {
		return err
	}

	outPath := filepath.Join(modelsDir, stem+".yaml")
	if _, err := os.Stat(outPath); err == nil && !overwrite {
		return fmt.Errorf("файл уже существует: %s (используйте --force)", outPath)
	}

	f, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(prepared); err != nil {
		return err
	}

	return enc.Close()
}
