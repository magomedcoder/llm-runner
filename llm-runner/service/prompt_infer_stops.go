package service

import "strings"

var (
	chatMLImStart = string([]byte{'<', '|', 'i', 'm', '_', 's', 't', 'a', 'r', 't', '|', '>'})
	chatMLImEnd   = string([]byte{'<', '|', 'i', 'm', '_', 'e', 'n', 'd', '|', '>'})
	chatMLImEndFW = "<\uff5cim_end\uff5c>"
	llamaEOT      = string([]byte{'<', '|', 'e', 'o', 't', '_', 'i', 'd', '|', '>'})
	llamaHdrStart = string([]byte{'<', '|', 's', 't', 'a', 'r', 't', '_', 'h', 'e', 'a', 'd', 'e', 'r', '_', 'i', 'd', '|', '>'})
	llamaHdrEnd   = string([]byte{'<', '|', 'e', 'n', 'd', '_', 'h', 'e', 'a', 'd', 'e', 'r', '_', 'i', 'd', '|', '>'})
)

func inferStopSequencesFromPrompt(prompt string) []string {
	seen := make(map[string]struct{})
	var out []string
	add := func(s string) {
		if s == "" {
			return
		}

		if _, ok := seen[s]; ok {
			return
		}

		seen[s] = struct{}{}
		out = append(out, s)
	}

	if strings.Contains(prompt, chatMLImStart) || strings.Contains(prompt, chatMLImEnd) {
		add(chatMLImStart)
		add(chatMLImEnd)
		add(chatMLImEndFW)
		add(chatMLImEnd + "\n")
		add(chatMLImEnd + "\r\n")
	}

	if strings.Contains(prompt, llamaEOT) || strings.Contains(prompt, llamaHdrStart) {
		add(llamaEOT)
		add(llamaEOT + "\n")
		add(llamaEOT + "\r\n")
		add(llamaHdrStart)
		add(llamaHdrEnd)
	}

	return out
}
