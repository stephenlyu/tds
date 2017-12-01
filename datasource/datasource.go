package datasource

import (
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
)

type InfoExDataSource interface {
	GetStockInfoEx(security *Security) (error, []InfoExItem)
	SetInfoEx(infoEx map[string][]InfoExItem) error
}

type StockNameItem struct {
	Date uint32
	Name string
}

type DataSource interface {
	InfoExDataSource

	// market - 市场代码, sz-深交所， sh-上交所
	GetStockCodes(exchange string) []string
	GetStockNameHistory(security *Security) []StockNameItem
	GetStockName(security *Security) string

	Reset()

	GetData(security *Security, period Period) (error, []Record)
	GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record)
	GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record)
	GetLastRecord(security *Security, period Period) (error, *Record)

	GetForwardAdjustedData(security *Security, period Period) (error, []Record)
	GetForwardAdjustedRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record)
	GetForwardAdjustedDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record)

	AppendData(security *Security, period Period, data []Record) error // Append data
	AppendRawData(security *Security, period Period, data []byte) error // Append raw data
	SaveData(security *Security, period Period, data []Record) error // Replace data with new data
}
