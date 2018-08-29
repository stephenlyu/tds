package util

import "math"

func Sum(values []float64) float64 {
	ret := 0.
	for _, v := range values {
		ret += v
	}
	return ret
}

func Mean(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	return Sum(values) / float64(len(values))
}

func Std(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	m := Mean(values)
	ss := 0.
	for _, v := range values {
		d := v - m
		ss += d * d
	}

	return math.Sqrt(ss / float64(len(values)))
}
