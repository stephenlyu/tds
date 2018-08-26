package stats

import (
	"math"
)

func CalcAPR(y []float64, scale float64) float64 {
	if len(y) == 0 {
		return 0
	}

	ret := y[len(y) - 1] / y[0]

	apr := math.Pow(ret, scale / float64(len(y))) - 1.0
	if math.IsNaN(apr) {
		apr = -100.0
	}
	return apr
}
