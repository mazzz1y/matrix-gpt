package text

import (
	"strings"
)

// HasPrefixIgnoreCase checks if string `s` has the given `prefix` irrespective of their case.
func HasPrefixIgnoreCase(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

// ReplaceIgnoreCase replaces occurrences of `old` in `s` with `new` irrespective of their case for a specified count `n`.
func ReplaceIgnoreCase(s, old, new string, n int) string {
	if old == "" {
		return s
	}

	msgLower := strings.ToLower(s)
	oldLower := strings.ToLower(old)
	var result strings.Builder
	count := 0

	for {
		if count == n && n >= 0 {
			result.WriteString(s)
			break
		}

		idx := strings.Index(msgLower, oldLower)
		if idx == -1 {
			result.WriteString(s)
			break
		}

		result.WriteString(s[:idx])
		result.WriteString(new)

		s = s[idx+len(old):]
		msgLower = msgLower[idx+len(old):]
		count++
	}

	return result.String()
}

// ReplaceAllIgnoreCase replaces all occurrences of `old` in `s` with `new` irrespective of their case.
func ReplaceAllIgnoreCase(s, old, new string) string {
	return ReplaceIgnoreCase(s, old, new, -1)
}
