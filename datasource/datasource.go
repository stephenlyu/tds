package datasource

import (
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
)

type InfoExDataSource interface {
	GetStockInfoEx(code string) (error, []InfoExItem)
	SetInfoEx(infoEx map[string][]InfoExItem) error
}

type DataSource interface {
	InfoExDataSource
	Reset()

	GetData(code string, period Period) (error, []Record)
	GetRangeData(code string, period Period, startDate, endDate uint64) (error, []Record)
	GetDataFromLast(code string, period Period, endDate uint64, count int) (error, []Record)

	GetForwardAdjustedData(code string, period Period) (error, []Record)
	GetForwardAdjustedRangeData(code string, period Period, startDate, endDate uint64) (error, []Record)
	GetForwardAdjustedDataFromLast(code string, period Period, endDate uint64, count int) (error, []Record)

	AppendData(code string, period Period, data []Record) error // Append data
	SaveData(code string, period Period, data []Record) error // Replace data with new data
}
