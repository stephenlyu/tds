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

type _IncrementStd struct {
	historyLen int
	historySum float64
	historyVar float64
}

func NewIncrementStd() *_IncrementStd {
	return &_IncrementStd{}
}

func (this *_IncrementStd) Feed(value float64) float64 {
	if this.historyLen == 0 {
		this.historyLen++
		this.historySum = value
		return this.historyVar
	}

	newSum := this.historySum + value

	M := float64(this.historyLen)
	const N = 1
	newLen := M + N

	newMean := newSum / newLen
	oldMean := this.historySum / M
	meanDelta := newMean - oldMean

	newVar := (M * (this.historyVar +  meanDelta * meanDelta) + (newMean - value) * (newMean - value)) / newLen

	this.historyLen++
	this.historySum = newSum
	this.historyVar = newVar

	return math.Sqrt(newVar)
}
