package activities

import (
	"regexp"
	"strings"
)

func sanitizeKey(key string) string {
	if key == "" {
		return ""
	}

	validKeyRegex := regexp.MustCompile(`[^_0-9A-Za-z]`)

	key = strings.TrimSpace(key)
	key = validKeyRegex.ReplaceAllString(key, "_")

	if key[0] >= '0' && key[0] <= '9' {
		key = "_" + key
	}

	return key
}

func sanitizeValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return sanitizeMap(val)
	case []any:
		out := make([]any, 0, len(val))
		for _, item := range val {
			out = append(out, sanitizeValue(item))
		}
		return out
	default:
		return val
	}
}

func sanitizeMap(m map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		out[sanitizeKey(k)] = sanitizeValue(v)
	}
	return out
}
