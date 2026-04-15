package mcpclient

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func WritePrometheusMetrics(w io.Writer) error {
	m := MCPCountersMap()
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		name := prometheusMetricName(k)
		fmt.Fprintf(w, "# TYPE %s counter\n", name)
		fmt.Fprintf(w, "%s %d\n", name, m[k])
	}

	writeCallToolDurationHistogram(w)

	return nil
}

func prometheusMetricName(key string) string {
	var b strings.Builder
	b.WriteString("gen_mcp_")
	for _, r := range key {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
		case r == '_':
			b.WriteRune('_')
		default:
			b.WriteRune('_')
		}
	}

	return b.String()
}
