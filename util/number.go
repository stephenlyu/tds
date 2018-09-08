package util

import (
	"strings"
	"fmt"
)

func TrimFloatStringZero (s string) string {
	i := strings.LastIndex(s, ".")
	if i < 0 {
		return s
	}

	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

func FormatFloat64(v float64) string {
	return TrimFloatStringZero(fmt.Sprintf("%f", v))
}
