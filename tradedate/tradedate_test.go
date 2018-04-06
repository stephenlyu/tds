package tradedate

import (
	"testing"
	"fmt"
)

func TestGetTradeDateRange(t *testing.T) {
	startTs, endTs := GetTradeDateRange(nil, "20180405")
	fmt.Println(startTs, endTs)
}