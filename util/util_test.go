package util

import (
	"testing"
	"fmt"
	"github.com/docker/docker/pkg/random"
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

func TestNewRingBuffer(t *testing.T) {
	rb := NewRingBuffer(10)
	for i := 0; i < 93; i++ {
		rb.Append(i)
	}
	fmt.Println(rb.Top)
	for i := 0; i < rb.Length; i++ {
		fmt.Println(rb.Get(i))
	}
}


func Test_IncrementStd_Feed(t *testing.T) {
	values := make([]float64, 100000)
	for i := range values {
		values[i] = random.Rand.Float64()
	}

	iStd := NewIncrementStd()

	var ret float64
	for _, v := range values {
		ret = iStd.Feed(v)
	}
	std := Std(values)
	fmt.Println(std, ret)
}