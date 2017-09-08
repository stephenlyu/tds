package datasource

import (
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
)

type InfoExDataSource interface {
	GetStockInfoEx(security *Security) (error, []InfoExItem)
	SetInfoEx(infoEx map[string][]InfoExItem) error
}

type DataSource interface {
	InfoExDataSource
	Reset()

	GetData(security *Security, period Period) (error, []Record)
	GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record)
	GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record)

	GetForwardAdjustedData(security *Security, period Period) (error, []Record)
	GetForwardAdjustedRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record)
	GetForwardAdjustedDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record)

	AppendData(security *Security, period Period, data []Record) error // Append data
	SaveData(security *Security, period Period, data []Record) error // Replace data with new data
}
