package stats

import "math"

func CalcMaxDrawDown(values []float64) (pos int, maxDrawDown float64) {
	peak := values[0]

	for i, v := range values {
		peak = math.Max(peak, v)
		v := (peak - v) / peak
		if v > maxDrawDown {
			maxDrawDown = v
			pos = i
		}
	}
	return
}
