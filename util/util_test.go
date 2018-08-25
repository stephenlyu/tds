package util

import (
	"testing"
	"fmt"
)

func TestUnzipFile(t *testing.T) {
	UnzipFile("zhb.zip", "temp")
}

func TestInSortedStrings(t *testing.T) {
	a := []string {"a", "b", "c", "d"}

	for _, v := range a {
		Assert(InSortedStrings(v, a), "")
	}

	Assert(!InSortedStrings(" ", a), "")
	Assert(!InSortedStrings("z", a), "")
}

func TestIterTools(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{3, 4, 5}
	c := []float64{5, 6, 7}

	r := Production(a, b, c)
	fmt.Println(r)
}