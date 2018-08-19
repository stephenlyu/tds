package util

import "sort"

func InStrings(v string, a []string) bool {
	for _, vv := range a {
		if v == vv {
			return true
		}
	}
	return false
}

func InInts(v int, a []int) bool {
	for _, vv := range a {
		if v == vv {
			return true
		}
	}
	return false
}

func InSortedStrings(v string, a []string) bool {
	index := sort.SearchStrings(a, v)
	if index == len(a) {
		return false
	}

	return a[index] == v
}
