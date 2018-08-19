package util

import "sort"

func SearchUInt64s(a []uint64, x uint64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func FindUInt64s(a []uint64, x uint64) int {
	index := SearchUInt64s(a, x)
	if index == len(a) {
		return -1
	}
	if a[index] != x {
		return -1
	}

	return index
}
