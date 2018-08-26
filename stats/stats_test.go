package stats

import (
	"testing"
	"fmt"
)

func TestSample(t *testing.T) {
	a := []float64{1,2,3,4,5,6,7}
	fmt.Println(Sample(a, 2))
}
