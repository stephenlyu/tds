package stats

import (
	"math"
	"github.com/stephenlyu/tds/util"
)

func CalcSharpeRatio(y []float64, scale float64) float64 {
	ret := make([]float64, len(y))
	ret[0] = 0

	for i := 0; i < len(y) - 1; i++ {
		ret[i + 1] = y[i + 1] / y[i] - 1
	}

	// 计算夏普比率
	mean := util.Mean(ret)
	std := util.Std(ret)
	var sharpe float64
	if math.IsNaN(mean) || math.IsNaN(std) || std == 0. {
		sharpe = -100.
	} else {
		sharpe = math.Sqrt(scale) * mean / std
	}

	return sharpe
}
