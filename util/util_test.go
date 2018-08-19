package util

import "testing"

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
