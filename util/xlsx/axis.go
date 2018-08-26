package xlsx

import (
	"strings"
	"github.com/stephenlyu/tds/util"
	"fmt"
)

var columnNames []string
var LETTERS = strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")

func init() {
	columnNames = append(columnNames, LETTERS...)

	comps := util.ProductionString(LETTERS, LETTERS)
	for _, a := range comps {
		columnNames = append(columnNames, strings.Join(a, ""))
	}

	comps = util.ProductionString(LETTERS, LETTERS, LETTERS)
	for _, a := range comps {
		columnNames = append(columnNames, strings.Join(a, ""))
	}
}

func GetAxis(row, col int) string {
	return fmt.Sprintf("%s%d", columnNames[col], row + 1)
}
