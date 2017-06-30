package tds

type InfoExDataSource interface {
	GetStockInfoEx(code string) (error, []InfoExItem)
	SetInfoEx(infoEx map[string][]InfoExItem) error
}

type DataSource interface {
	InfoExDataSource
	Reset()

	GetData(code string, period Period) (error, []Record)
	GetRangeData(code string, period Period, startDate, endDate Date) (error, []Record)
	GetDataFromLast(code string, period Period, endDate Date, count int) (error, []Record)

	GetForwardAdjustedData(code string, period Period) (error, []Record)
	GetForwardAdjustedRangeData(code string, period Period, startDate, endDate Date) (error, []Record)
	GetForwardAdjustedDataFromLast(code string, period Period, endDate Date, count int) (error, []Record)
}

type datasource struct {
	Root string

	InfoEx map[string][]InfoExItem
}

func NewDataSource(dsDir string) DataSource {
	return &datasource{Root: dsDir}
}

func (this *datasource) Reset() {
	this.InfoEx = nil
}

func (this *datasource) GetStockInfoEx(code string) (error, []InfoExItem){
	return nil, nil
}

func (this *datasource) SetInfoEx(infoEx map[string][]InfoExItem) error {
	return nil
}

func (this *datasource) GetData(code string, period Period) (error, []Record) {
	return nil, nil
}

func (this *datasource) GetRangeData(code string, period Period, startDate, endDate Date) (error, []Record) {
	return nil, nil
}

func (this *datasource) GetDataFromLast(code string, period Period, endDate Date, count int) (error, []Record) {
	return nil, nil
}

func (this *datasource) GetForwardAdjustedData(code string, period Period) (error, []Record) {
	return nil, nil
}

func (this *datasource) GetForwardAdjustedRangeData(code string, period Period, startDate, endDate Date) (error, []Record) {
	return nil, nil
}

func (this *datasource) GetForwardAdjustedDataFromLast(code string, period Period, endDate Date, count int) (error, []Record) {
	return nil, nil
}
