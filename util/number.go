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

func MaxUInt64(a, b uint64) uint64 {
	if a < b {
		return b
	}
	return a
}

func MinUInt64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
