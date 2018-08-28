package util

import "strings"

func TrimFloatStringZero (s string) string {
	i := strings.LastIndex(s, ".")
	if i < 0 {
		return s
	}

	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}
