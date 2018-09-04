package stats

import (
	"github.com/stephenlyu/tds/util"
)

func Sample(y []float64, rate int) []float64 {
	util.Assert(rate > 0, "")

	size := len(y) / rate
	ret := make([]float64, size)

	for i, j := rate - 1, 0; i < len(y); i, j = i + rate, j + 1 {
		ret[j] = y[i]
	}
	return ret
}

func SampleString(y []string, rate int) []string {
	util.Assert(rate > 0, "")

	size := len(y) / rate
	ret := make([]string, size)

	for i, j := rate - 1, 0; i < len(y); i, j = i + rate, j + 1 {
		ret[j] = y[i]
	}
	return ret
}
