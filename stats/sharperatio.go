package stats

import "math"

func sum(values []float64) float64 {
	ret := 0.
	for _, v := range values {
		ret += v
	}
	return ret
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	return sum(values) / float64(len(values))
}

func std(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	m := mean(values)
	ss := 0.
	for _, v := range values {
		d := v - m
		ss += d * d
	}

	return math.Sqrt(ss / float64(len(values)))
}

func CalcSharpeRatio(y []float64, scale float64) float64 {
	ret := make([]float64, len(y))
	ret[0] = 0

	for i := 0; i < len(y) - 1; i++ {
		ret[i + 1] = y[i + 1] / y[i] - 1
	}

	// 计算夏普比率
	mean := mean(ret)
	std := std(ret)
	var sharpe float64
	if math.IsNaN(mean) || math.IsNaN(std) || std == 0. {
		sharpe = -100.
	} else {
		sharpe = math.Sqrt(scale) * mean / std
	}

	return sharpe
}
