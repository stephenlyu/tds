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

// endDate: inclusive
type BaseDataSource interface {
	GetData(security *Security, period Period) (error, []Record)
	GetDataEx(security *Security, period Period, startDate uint64, count int) (error, []Record)

	GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record)
	GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record)
	GetLastRecord(security *Security, period Period) (error, *Record)

	GetForwardAdjustedRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record)

	AppendData(security *Security, period Period, data []Record) error // Append data
	SaveData(security *Security, period Period, data []Record) error   // Replace data with new data

	RemoveData(security *Security, period Period, startDate, endDate uint64) error
}

type CompositeDataSource interface {
	BaseDataSource
	AddSubDatasource(ds BaseDataSource)
}

type DateRange struct {
	Security  *Security
	StartDate uint64 // inclusive, 0 if no start date
	EndDate   uint64 // exclusive, 0 if no end date
}

type DateRangeMapper interface {
	MapDateRanges(security *Security) []DateRange
}

type MappedDataSource interface {
	BaseDataSource
	SetMapper(mapper DateRangeMapper)
	SetTargetDataSource(ds BaseDataSource)
}

type DataSource interface {
	InfoExDataSource
	BaseDataSource

	// market - 市场代码, sz-深交所， sh-上交所
	GetStockCodes(exchange string) []string
	GetStockNameHistory(security *Security) []StockNameItem
	GetStockName(security *Security) string
	GetStockNames() map[string]string

	Reset()

	SupportedPeriods() []Period

	GetForwardAdjustedData(security *Security, period Period) (error, []Record)
	GetForwardAdjustedRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record)
	GetForwardAdjustedDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record)

	AppendRawData(security *Security, period Period, data []byte) error // Append raw data

	// Remove data which date is greater or equal to date
	TruncateTo(security *Security, period Period, date uint64) error
}

type TickDataSource interface {
	Save(security *Security, ticks []TickItem) error
	Get(security *Security, startTs, endTs uint64) (error, []TickItem)
	Remove(security *Security, startTs, endTs uint64) error
}
