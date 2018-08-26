package xlsx

import (
	"testing"
	"github.com/stephenlyu/tds/util"
	"fmt"
)

func TestWrite(t *testing.T) {
	records := [][]interface{} {
		[]interface{}{1, 3.0, "first", 9.8, 255, 888},
		[]interface{}{2, 8.0, "second\"", 19.8, 2555, 999},
	}
	err := Write("test.xlsx", records)
	util.Assert(err == nil, fmt.Sprintf("%+v", err))
}
