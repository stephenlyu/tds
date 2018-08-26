package util

func Production(args ...[]float64) [][]float64 {
	var ret [][]float64 = [][]float64{[]float64{}}
	for _, arg := range args {
		var newRet [][]float64
		for _, values := range ret {
			for _, v := range arg {
				newValues := make([]float64, len(values) + 1)
				copy(newValues[:len(values)], values)
				newValues[len(values)] = v
				newRet = append(newRet, newValues)
			}
		}
		ret = newRet
	}
	return ret
}

func ProductionString(args ...[]string) [][]string {
	var ret [][]string = [][]string{[]string{}}
	for _, arg := range args {
		var newRet [][]string
		for _, values := range ret {
			for _, v := range arg {
				newValues := make([]string, len(values) + 1)
				copy(newValues[:len(values)], values)
				newValues[len(values)] = v
				newRet = append(newRet, newValues)
			}
		}
		ret = newRet
	}
	return ret
}
