package tdxdatasource

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"sort"
	"errors"

	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
	. "github.com/stephenlyu/tds/datasource"
	"encoding/json"
	"strings"
	"fmt"
	"github.com/z-ray/log"
	"regexp"
	"sync"
	"golang.org/x/text/encoding/simplifiedchinese"
	"encoding/binary"
)

var (
	linefeedPattern, _ = regexp.Compile("(\r\n|\r|\n)")
)

var periodNameMap = map[string]string {
	"fzline": "MINUTE5",
	"lday": "DAY1",
	"minline": "MINUTE1",
}

var periodNameReMap = map[string]string {
	"MINUTE5": "fzline",
	"DAY1": "lday",
	"MINUTE1": "minline",
}

var fileNameSuffixMap = map[string]string {
	"MINUTE1": ".lc1",
	"MINUTE5": ".lc5",
	"DAY1": ".day",
}

var exchangeBlockMap = map[string]string {
	"SZ": "0",
	"SH": "1",
}
var blockExchangeMap = map[string]string {
	"0": "SZ",
	"1": "SH",
}


type tdxDataSource struct {
	DataDir          string
	ConfigDir        string

	NeedBuildCache   bool

	InfoEx           map[string][]InfoExItem

	lock             sync.Mutex
	stockCodeCache   map[string][]string
	stockNameHistory map[string][]StockNameItem
	stockNames		 map[string]string
}

func NewDataSource(dsDir string, needBuildCache bool) DataSource {
	return &tdxDataSource{
		DataDir: filepath.Join(dsDir, "vipdoc"),
		ConfigDir: filepath.Join(dsDir, "T0002"),
		NeedBuildCache: needBuildCache,

		stockCodeCache: make(map[string][]string),
	}
}

func (this *tdxDataSource) Reset() {
	this.InfoEx = nil
}

func (this *tdxDataSource) GetStockCodes(exchange string) []string {
	exchange = strings.ToUpper(exchange)

	block, ok := exchangeBlockMap[exchange]
	if !ok {
		return nil
	}

	this.lock.Lock()
	defer this.lock.Unlock()
	if ret, ok := this.stockCodeCache[block]; ok {
		return ret
	}

	filePath := filepath.Join(this.ConfigDir, "hq_cache/tipinfo.dat")
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("tdxDataSource.GetStockCodes read file fail, error: %+v", err)
	}

	bytes = linefeedPattern.ReplaceAll(bytes, []byte("\n"))
	lines := strings.Split(string(bytes), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		block := parts[0]
		if _, ok := blockExchangeMap[block]; !ok {
			continue
		}

		code := parts[1]

		this.stockCodeCache[block] = append(this.stockCodeCache[block], fmt.Sprintf("%s.%s", code, blockExchangeMap[block]))
	}

	return this.stockCodeCache[block]
}

func (this *tdxDataSource) GetStockNameHistory(security *Security) []StockNameItem {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.stockNameHistory != nil {
		return this.stockNameHistory[security.Code]
	}

	this.stockNameHistory = make(map[string][]StockNameItem)

	// Load stock names
	filePath := filepath.Join(this.ConfigDir, "hq_cache/profile.dat")
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("tdxDataSource.GetStockNameHistory read file fail, error: %+v", err)
		return nil
	}

	elemSize := 64

	end := len(bytes) / elemSize * elemSize

	fmt.Println(len(bytes) / elemSize)

	gbkDecoder := simplifiedchinese.GBK.NewDecoder()

	for i := 0; i < end; i += elemSize {
		r := bytes[i:i+elemSize]

		code := string(r[1:7])
		nameGBKBytes := r[8:17]

		for j := len(nameGBKBytes) - 1; j >= 0; j-- {
			if nameGBKBytes[j] != 0 {
				nameGBKBytes = nameGBKBytes[:j+1]
				break
			}
		}

		nameBytes := make([]byte, 30)

		nDest, _, err := gbkDecoder.Transform(nameBytes, nameGBKBytes, true)
		if err != nil {
			log.Errorf("tdxDataSource.GetStockNameHistory, decode name fail, error: %v", err)
			continue
		}

		name := string(nameBytes[:nDest])
		date := binary.LittleEndian.Uint32(r[17:21])

		this.stockNameHistory[code] = append(this.stockNameHistory[code], StockNameItem{date, name})
	}

	for _, items := range this.stockNameHistory {
		sort.SliceStable(items, func (i, j int) bool {
			return items[i].Date > items[j].Date
		})
	}

	return this.stockNameHistory[security.Code]
}

func (this *tdxDataSource) GetStockName(security *Security) string {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.stockNames != nil {
		return this.stockNames[security.Code]
	}

	this.stockNames = make(map[string]string)

	// Load stock names
	filePath := filepath.Join(this.ConfigDir, "hq_cache/names.dat")
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("tdxDataSource.GetStockName read file fail, error: %+v", err)
		return ""
	}

	elemSize := 29

	end := len(bytes) / elemSize * elemSize

	gbkDecoder := simplifiedchinese.GBK.NewDecoder()

	for i := 0; i < end; i += elemSize {
		r := bytes[i:i+elemSize]

		code := string(r[0:6])
		var nameGBKBytes []byte
		for j := 8; j < 16; j++ {
			if r[j] == 0 {
				nameGBKBytes = r[8:j]
				break
			}
		}
		if nameGBKBytes == nil {
			nameGBKBytes = r[8:16]
		}

		nameBytes := make([]byte, 30)

		nDest, _, err := gbkDecoder.Transform(nameBytes, nameGBKBytes, true)
		if err != nil {
			log.Errorf("tdxDataSource.GetStockName, decode name fail, error: %v", err)
			continue
		}

		name := string(nameBytes[:nDest])
		this.stockNames[code] = name
	}

	return this.stockNames[security.Code]
}

func (this *tdxDataSource) GetStockInfoEx(security *Security) (error, []InfoExItem){
	if this.InfoEx == nil {
		filePath := filepath.Join(this.DataDir, "infoex.dat")

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err, nil
		}

		err = json.Unmarshal(bytes, &this.InfoEx)
		if err != nil {
			return err, nil
		}
	}
	code := SecurityToString(security)
	return nil, this.InfoEx[code]
}

func (this *tdxDataSource) SetInfoEx(infoEx map[string][]InfoExItem) error {
	this.InfoEx = infoEx
	filePath := filepath.Join(this.DataDir, "vipdoc/infoex.dat")

	bytes, err := json.Marshal(this.InfoEx)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, bytes, 0666)
}

func (this *tdxDataSource) GetData(security *Security, period Period) (error, []Record) {
	return this.GetRangeData(security, period, 0, 0)
}

func (this *tdxDataSource) getDataFile(security *Security, period Period) (Period, string) {
	code := SecurityToString(security)

	root := filepath.Join(this.DataDir, strings.ToLower(security.Exchange))

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Errorf("tdxDataSource.getDataFile read dir %s fail, error: %v", root, err)
		return nil, ""
	}

	periods := make([]Period, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		name, ok := periodNameMap[f.Name()]
		if !ok {
			continue
		}

		err, p := PeriodFromString(name)
		if err != nil {
			log.Debugf("tdxDataSource.getDataFile parse period %s error: %v", name, err)
			continue
		}
		if !p.CanConvertTo(period) {
			continue
		}

		filePath := filepath.Join(root, f.Name(), fmt.Sprintf("%s%s", code, fileNameSuffixMap[p.Name()]))
		_, err = os.Stat(filePath)
		if err != nil {
			log.Debugf("tdxDataSource.getDataFile stat file: %s error: %v", filePath, err)
			continue
		}

		periods = append(periods, p)
	}

	if len(periods) == 0 {
		log.Errorf("tdxDataSource.getDataFile no period directory found")
		return nil, ""
	}

	sort.SliceStable(periods, func (i,j int) bool {
		return periods[i].Gt(periods[j])
	})
	dataPeriod := periods[0]
	return dataPeriod, filepath.Join(root, periodNameReMap[dataPeriod.Name()], fmt.Sprintf("%s%s", code, fileNameSuffixMap[dataPeriod.Name()]))
}

func (this *tdxDataSource) binarySearchRecord(reader RecordReader, period Period, date uint64, count int) (error, int, bool) {
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

		if records[0].Date == date {
			return nil, mid, true
		} else if records[0].Date < date {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return nil, low, false
}

func (this *tdxDataSource) GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record) {
	if startDate != 0 && endDate != 0 && startDate > endDate {
		return nil, []Record{}
	}

	dataPeriod, dataFile := this.getDataFile(security, period)
	if dataFile == "" {
		return errors.New("data file not found"), nil
	}

	file, err := os.Open(dataFile)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	marshaller := NewMarshaller(period)

	reader := NewRecordReader(file, TDX_RECORD_SIZSE, marshaller)
	err, recordCount := reader.Count()
	if err != nil {
		return err, nil
	}

	var startIndex = 0
	var endIndex = recordCount
	if startDate != 0 {
		err, startIndex, _ = this.binarySearchRecord(reader, dataPeriod, startDate, recordCount)
		if err != nil {
			return err, nil
		}
	}
	if endDate != 0 {
		var found bool
		err, endIndex, found = this.binarySearchRecord(reader, dataPeriod, endDate, recordCount)
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

func (this *tdxDataSource) GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record) {
	// FIXME: 获取其他周期数据时，如何保证count?
	dataPeriod, dataFile := this.getDataFile(security, period)
	if dataFile == "" {
		return errors.New("data file not found"), nil
	}

	file, err := os.Open(dataFile)
	if err != nil {
		return err, nil
	}
	defer file.Close()

	marshaller := NewMarshaller(period)

	reader := NewRecordReader(file, TDX_RECORD_SIZSE, marshaller)
	err, recordCount := reader.Count()
	if err != nil {
		return err, nil
	}

	var endIndex = recordCount
	if endDate != 0 {
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

func (this *tdxDataSource) GetForwardAdjustedData(security *Security, period Period) (error, []Record) {
	return this.GetForwardAdjustedRangeData(security, period, 0, 0)
}

func (this *tdxDataSource) GetForwardAdjustedRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record) {
	err, records := this.GetRangeData(security, period, startDate, endDate)
	if err != nil {
		return err, nil
	}

	err, exItems := this.GetStockInfoEx(security)
	if err != nil {
		return err, nil
	}
	if len(exItems) == 0 {
		return err, records
	}

	converter := NewForwardAdjustConverter(period, exItems)
	return nil, converter.Convert(records)
}

func (this *tdxDataSource) GetForwardAdjustedDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record) {
	err, records := this.GetDataFromLast(security, period, endDate, count)
	if err != nil {
		return err, nil
	}

	err, exItems := this.GetStockInfoEx(security)
	if err != nil {
		return err, nil
	}
	if len(exItems) == 0 {
		return err, records
	}

	converter := NewForwardAdjustConverter(period, exItems)
	return nil, converter.Convert(records)
}

func (this *tdxDataSource) checkData(period Period, data []Record) bool {
	for i := 0; i < len(data) - 1; i++ {
		if data[i].Date >= data[i + 1].Date {
			return false
		}
	}
	return true
}

func (this *tdxDataSource) AppendData(security *Security, period Period, data []Record) error {
	if len(data) == 0 {
		return nil
	}

	if !this.checkData(period, data) {
		return errors.New("bad data")
	}

	code := SecurityToString(security)
	filePath := filepath.Join(this.DataDir, period.ShortName(), code)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	marshaller := NewMarshaller(period)
	reader := NewRecordReader(file, TDX_RECORD_SIZSE, marshaller)
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
		lastDate := records[0].Date

		for i := 0; i < len(data); i++ {
			r := data[i]
			if r.Date > lastDate {
				data = data[i:]
				break
			}
		}
	}

	writer := NewRecordWriter(file, TDX_RECORD_SIZSE, marshaller)
	return writer.Write(recordCount, data)
}

func (this *tdxDataSource) SaveData(security *Security, period Period, data []Record) error {
	if len(data) == 0 {
		return nil
	}

	if !this.checkData(period, data) {
		return errors.New("bad data")
	}

	filePath := filepath.Join(this.DataDir, period.ShortName(), SecurityToString(security))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	marshaller := NewMarshaller(period)
	writer := NewRecordWriter(file, TDX_RECORD_SIZSE, marshaller)
	return writer.Write(0, data)
}
