package tds

import (
	"encoding/gob"
	"path/filepath"
	"os"
	"io/ioutil"
	"sort"
	"errors"
)

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

	AppendData(code string, period Period, data []Record) error // Append data
	SaveData(code string, period Period, data []Record) error // Replace data with new data
}

type datasource struct {
	Root string
	NeedBuildCache bool

	InfoEx map[string][]InfoExItem
}

func NewDataSource(dsDir string, needBuildCache bool) DataSource {
	return &datasource{Root: dsDir, NeedBuildCache: needBuildCache}
}

func (this *datasource) Reset() {
	this.InfoEx = nil
}

func (this *datasource) GetStockInfoEx(code string) (error, []InfoExItem){
	if this.InfoEx == nil {
		filePath := filepath.Join(this.Root, "infoex.dat")
		file, err := os.Open(filePath)
		if err != nil {
			return err, nil
		}
		defer file.Close()
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&this.InfoEx)
		if err != nil {
			return err, nil
		}
	}

	return nil, this.InfoEx[code]
}

func (this *datasource) SetInfoEx(infoEx map[string][]InfoExItem) error {
	this.InfoEx = infoEx
	filePath := filepath.Join(this.Root, "infoex.dat")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(infoEx)
}

func (this *datasource) GetData(code string, period Period) (error, []Record) {
	return this.GetRangeData(code, period, nil, nil)
}

func (this *datasource) getDataFile(code string, period Period) (Period, string) {
	files, err := ioutil.ReadDir(this.Root)
	if err != nil {
		return nil, ""
	}

	periods := make([]Period, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		err, p := FromString(f.Name())
		if err != nil {
			continue
		}
		if !p.CanConvertTo(period) {
			continue
		}

		filePath := filepath.Join(this.Root, f.Name(), code)
		_, err = os.Stat(filePath)
		if err != nil {
			continue
		}

		periods = append(periods, p)
	}

	if len(periods) == 0 {
		return nil, ""
	}

	sort.SliceStable(periods, func (i,j int) bool {
		return periods[i].Gt(periods[j])
	})

	return periods[0], filepath.Join(this.Root, periods[0].ShortName(), code)
}

func (this *datasource) binarySearchRecord(reader RecordReader, period Period, date Date, count int) (error, int, bool) {
	low := 0
	high := count - 1
	var mid int

	for low <= high {
		mid = (low + high) / 2
		err, records := reader.Read(mid, mid + 1)
		if err != nil {
			return err, -1, false
		}
		if len(records) == 0 {
			return errors.New("no data read"), -1, false
		}

		rDate := NewPeriodDate(period, records[0].Date)
		if rDate.Eq(date) {
			return nil, mid, true
		} else if rDate.Lt(date) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return nil, low, false
}

func (this *datasource) GetRangeData(code string, period Period, startDate, endDate Date) (error, []Record) {
	if startDate != nil && endDate != nil && startDate.Gt(endDate) {
		return nil, []Record{}
	}

	dataPeriod, dataFile := this.getDataFile(code, period)
	if dataFile == "" {
		return errors.New("data file not found"), nil
	}

	file, err := os.Open(dataFile)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	reader := NewRecordReader(file)
	err, recordCount := reader.Count()
	if err != nil {
		return err, nil
	}

	var startIndex = 0
	var endIndex = recordCount
	if startDate != nil {
		err, startIndex, _ = this.binarySearchRecord(reader, dataPeriod, startDate, recordCount)
		if err != nil {
			return err, nil
		}
	}
	if endDate != nil {
		err, endIndex, found := this.binarySearchRecord(reader, dataPeriod, endDate, recordCount)
		if err != nil {
			return err, nil
		}
		if found {
			endIndex++
		}
	}

	err, records := reader.Read(startIndex, endIndex)
	if err != nil {
		return err, nil
	}

	if period.Eq(dataPeriod) {
		return nil, records
	}

	converter := NewPeriodConverter(dataPeriod, period)
	return nil, converter.Convert(records)
}

func (this *datasource) GetDataFromLast(code string, period Period, endDate Date, count int) (error, []Record) {
	dataPeriod, dataFile := this.getDataFile(code, period)
	if dataFile == "" {
		return errors.New("data file not found"), nil
	}

	file, err := os.Open(dataFile)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	reader := NewRecordReader(file)
	err, recordCount := reader.Count()
	if err != nil {
		return err, nil
	}

	var endIndex = recordCount
	if endDate != nil {
		err, endIndex, found := this.binarySearchRecord(reader, dataPeriod, endDate, recordCount)
		if err != nil {
			return err, nil
		}
		if found {
			endIndex++
		}
	}

	startIndex := endIndex - count
	if startIndex < 0 {
		startIndex = 0
	}

	err, records := reader.Read(startIndex, endIndex)
	if err != nil {
		return err, nil
	}

	if period.Eq(dataPeriod) {
		return nil, records
	}

	converter := NewPeriodConverter(dataPeriod, period)
	return nil, converter.Convert(records)
}

func (this *datasource) GetForwardAdjustedData(code string, period Period) (error, []Record) {
	return this.GetForwardAdjustedRangeData(code, period, nil, nil)
}

func (this *datasource) GetForwardAdjustedRangeData(code string, period Period, startDate, endDate Date) (error, []Record) {
	err, records := this.GetRangeData(code, period, startDate, endDate)
	if err != nil {
		return err, nil
	}

	err, exItems := this.GetStockInfoEx(code)
	if err != nil {
		return err, nil
	}
	if len(exItems) == 0 {
		return err, records
	}

	converter := NewForwardAdjustConverter(period, exItems)
	return nil, converter.Convert(records)
}

func (this *datasource) GetForwardAdjustedDataFromLast(code string, period Period, endDate Date, count int) (error, []Record) {
	err, records := this.GetDataFromLast(code, period, endDate, count)
	if err != nil {
		return err, nil
	}

	err, exItems := this.GetStockInfoEx(code)
	if err != nil {
		return err, nil
	}
	if len(exItems) == 0 {
		return err, records
	}

	converter := NewForwardAdjustConverter(period, exItems)
	return nil, converter.Convert(records)
}

func (this *datasource) checkData(period Period, data []Record) bool {
	for i := 0; i < len(data) - 1; i++ {
		date := NewPeriodDate(period, data[i].Date)
		nextDate := NewPeriodDate(period, data[i + 1].Date)
		if !date.Lt(nextDate) {
			return false
		}
	}
	return true
}

func (this *datasource) AppendData(code string, period Period, data []Record) error {
	if len(data) == 0 {
		return nil
	}

	if !this.checkData(period, data) {
		return errors.New("bad data")
	}

	filePath := filepath.Join(this.Root, period.ShortName(), code)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := NewRecordReader(file)
	err, recordCount := reader.Count()
	if err != nil {
		return err
	}
	if recordCount > 0 {
		err, records := reader.Read(recordCount - 1, recordCount)
		if err != nil {
			return err
		}
		if len(records) == 0 {
			return errors.New("no data read")
		}
		lastDate := NewPeriodDate(period, records[0].Date)

		for i := 0; i < len(data); i++ {
			r := data[i]
			date := NewPeriodDate(period, r.Date)
			if date.Gt(lastDate) {
				data = data[i:]
				break
			}
		}
	}

	writer := NewRecordWriter(file)
	return writer.Write(recordCount, data)
}

func (this *datasource) SaveData(code string, period Period, data []Record) error {
	if len(data) == 0 {
		return nil
	}

	if !this.checkData(period, data) {
		return errors.New("bad data")
	}

	filePath := filepath.Join(this.Root, period.ShortName(), code)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := NewRecordWriter(file)
	return writer.Write(0, data)
}
