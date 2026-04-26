package handler

import "strings"

const MCPSecretMaskedPlaceholder = "***"

func maskSecretMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return map[string]string{}
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		if strings.TrimSpace(v) == "" {
			out[k] = v
		} else {
			out[k] = MCPSecretMaskedPlaceholder
		}
	}

	return out
}

func mergeMaskedSecretMaps(incoming, existing map[string]string) map[string]string {
	if incoming == nil {
		incoming = map[string]string{}
	}

	if existing == nil {
		existing = map[string]string{}
	}

	out := make(map[string]string, len(incoming))
	for k, v := range incoming {
		if strings.TrimSpace(v) == MCPSecretMaskedPlaceholder {
			if old, ok := existing[k]; ok {
				out[k] = old
			}

			continue
		}
		out[k] = v
	}

	return out
}
