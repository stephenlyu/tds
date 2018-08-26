package csvdatasource

import (
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/storage"
	"fmt"
	"path/filepath"
	"reflect"
	"os"
)

type _CSVDataSource struct {
	engine *storage.CsvEngine
	csvDir string
}

func NewCSVDataSource(csvDir string) datasource.BaseDataSource {
	return &_CSVDataSource{
		csvDir: csvDir,
		engine: storage.NewCsvEngine(reflect.TypeOf(Record{})),
	}
}

func (this *_CSVDataSource) filePath(security *Security, period Period) string {
	fileName := fmt.Sprintf("%s.%s.csv", period.ShortName(), security.String())
	return filepath.Join(this.csvDir, fileName)
}

func (this *_CSVDataSource) GetData(security *Security, period Period) (error, []Record) {
	filePath := this.filePath(security, period)

	err, data := this.engine.Load(filePath)
	if err != nil {
		return err, nil
	}

	ret := make([]Record, len(data))
	for i := range data {
		r := data[i].(*Record)
		ret[i] = *r
	}
	return nil, ret
}

func (this *_CSVDataSource) GetDataEx(security *Security, period Period, startDate uint64, count int) (error, []Record) {
	err, records := this.GetData(security, period)
	if err != nil {
		return err, nil
	}

	for i := range records {
		if records[i].Date >= startDate {
			end := i + count
			if end > len(records) {
				end = len(records)
			}
			ret := make([]Record, end - i)
			copy(ret, records[i:end])
			return nil, ret
		}
	}
	return nil, nil
}

func (this *_CSVDataSource) GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record) {
	err, records := this.GetData(security, period)
	if err != nil {
		return err, nil
	}

	var start, end int
	if startDate > 0 {
		for i := range records {
			if records[i].Date >= startDate {
				start = i
				break
			}
		}
	}
	end = len(records)
	if endDate > 0 {
		for i := start + 1; i < len(records); i++ {
			if records[i].Date > endDate {
				end = i
				break
			}
		}
	}

	if end > start {
		ret := make([]Record, end - start)
		copy(ret, records[start:end])
		return nil, ret
	}

	return nil, nil
}

func (this *_CSVDataSource) GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record) {
	panic("unimplented")
	return nil, nil
}

func (this *_CSVDataSource) GetLastRecord(security *Security, period Period) (error, *Record) {
	panic("unimplented")
	return nil, nil
}

func (this *_CSVDataSource) AppendData(security *Security, period Period, data []Record) error {
	panic("unimplented")
	return nil
}

func (this *_CSVDataSource) SaveData(security *Security, period Period, data []Record) error {
	os.MkdirAll(this.csvDir, 0777)
	filePath := this.filePath(security, period)
	dataCpy := make([]interface{}, len(data))
	for i := range data {
		dataCpy[i] = &data[i]
	}

	return this.engine.Save(filePath, dataCpy)
}

func (this *_CSVDataSource) RemoveData(security *Security, period Period, startDate, endDate uint64) error {
	panic("unimplented")
	return nil
}
