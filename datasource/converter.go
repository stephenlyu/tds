package datasource

import (
	"sort"
	"errors"

	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
	. "github.com/stephenlyu/tds/date"
)

// TODO: 数据转换时，针对跨市场的情况，需要重新考虑时间的计算问题

type Converter interface {
	Convert(sourceData []Record) []Record
}

type periodConverter struct {
	srcPeriod Period
	destPeriod Period
}

type forwardAdjustConverter struct {
	period Period
	items []InfoExItem
}

// Period Data Converters

func NewPeriodConverter(srcPeriod Period, destPeriod Period) Converter {
	if !srcPeriod.CanConvertTo(destPeriod) {
		return nil
	}

	return &periodConverter{srcPeriod: srcPeriod, destPeriod: destPeriod}
}

func MergeData(destData []Record, sourceData []Record, multiplier int) []Record {
	for i := 0; i < len(sourceData); i++ {
		r := &sourceData[i]
		if i % multiplier == 0 {
			destData = append(destData, *r)
			continue
		}

		dr := &destData[len(destData) - 1]

		dr.Close = r.Close
		if r.Low < dr.Low {
			dr.Low = r.Low
		}
		if r.High > dr.High {
			dr.High = r.High
		}
		dr.Date = r.Date

		dr.Volume += r.Volume
		dr.Amount += r.Amount
	}
	return destData
}


func (this *periodConverter) doMerge(destData []Record, sourceData []Record, multiplier int) []Record {
	return MergeData(destData, sourceData, multiplier)
}

func (this *periodConverter) convertMinute2Minute(sourceData []Record) []Record {
	multiplier := this.destPeriod.UnitCount() / this.srcPeriod.UnitCount()

	destData := make([]Record, 0, len(sourceData))

	var lastDay int = 0

	startIndex := -1
	for i := 0; i < len(sourceData); i++ {
		r := &sourceData[i]
		day := GetDateDay(r.Date)
		if day != lastDay {
			if startIndex >= 0 {
				destData = this.doMerge(destData, sourceData[startIndex:i], multiplier)
			}
			startIndex = i
		}

		lastDay = day
		if i == len(sourceData) - 1 {
			destData = this.doMerge(destData, sourceData[startIndex:], multiplier)
		}
	}

	ret := make([]Record, len(destData))
	copy(ret, destData)
	return ret
}

func (this *periodConverter) convertMinute2Day(sourceData []Record) []Record {
	destData := make([]Record, 0, len(sourceData))

	var lastDay uint64 = 0

	startIndex := -1
	for i := 0; i < len(sourceData); i++ {
		r := &sourceData[i]
		day := GetDayTimestamp(r.Date)
		if day != lastDay{
			if startIndex >= 0 {
				destData = this.doMerge(destData, sourceData[startIndex:i], i - startIndex)
				destData[len(destData) - 1].Date = lastDay
			}
			startIndex = i
		}

		lastDay = day
		if i == len(sourceData) - 1 {
			destData = this.doMerge(destData, sourceData[startIndex:], i + 1 - startIndex)
			destData[len(destData) - 1].Date = lastDay
		}
	}

	ret := make([]Record, len(destData))
	copy(ret, destData)
	return ret
}

func (this *periodConverter) convertDay2Week(sourceData []Record) []Record {
	destData := make([]Record, 0, len(sourceData))

	var lastWeek int = 0

	startIndex := -1
	for i := 0; i < len(sourceData); i++ {
		r := &sourceData[i]
		week := GetDateWeek(r.Date)
		if week != lastWeek {
			if startIndex >= 0 {
				destData = this.doMerge(destData, sourceData[startIndex:i], i - startIndex)
			}
			startIndex = i
		}

		lastWeek = week
		if i == len(sourceData) - 1 {
			destData = this.doMerge(destData, sourceData[startIndex:], i + 1 - startIndex)
		}
	}

	ret := make([]Record, len(destData))
	copy(ret, destData)
	return ret
}

func (this *periodConverter) convertDay2Month(sourceData []Record) []Record {
	destData := make([]Record, 0, len(sourceData))

	var lastMonth int = 0

	startIndex := -1
	for i := 0; i < len(sourceData); i++ {
		r := &sourceData[i]
		month := GetDateMonth(r.Date)
		if month != lastMonth {
			if startIndex >= 0 {
				destData = this.doMerge(destData, sourceData[startIndex:i], i - startIndex)
			}
			startIndex = i
		}

		lastMonth = month
		if i == len(sourceData) - 1 {
			destData = this.doMerge(destData, sourceData[startIndex:], i + 1 - startIndex)
		}
	}

	ret := make([]Record, len(destData))
	copy(ret, destData)
	return ret
}

func (this *periodConverter) convert2Quarter(sourceData []Record) []Record {
	destData := make([]Record, 0, len(sourceData))

	var lastQuarter int = 0

	startIndex := -1
	for i := 0; i < len(sourceData); i++ {
		r := &sourceData[i]
		quarter := GetDateQuarter(r.Date)
		if quarter != lastQuarter {
			if startIndex >= 0 {
				destData = this.doMerge(destData, sourceData[startIndex:i], i - startIndex)
			}
			startIndex = i
		}

		lastQuarter = quarter
		if i == len(sourceData) - 1 {
			destData = this.doMerge(destData, sourceData[startIndex:], i + 1 - startIndex)
		}
	}

	ret := make([]Record, len(destData))
	copy(ret, destData)
	return ret
}

func (this *periodConverter) convert2Year(sourceData []Record) []Record {
	destData := make([]Record, 0, len(sourceData))

	var lastYear int = 0

	startIndex := -1
	for i := 0; i < len(sourceData); i++ {
		r := &sourceData[i]
		year := GetDateYear(r.Date)
		if year != lastYear {
			if startIndex >= 0 {
				destData = this.doMerge(destData, sourceData[startIndex:i], i - startIndex)
			}
			startIndex = i
		}

		lastYear = year
		if i == len(sourceData) - 1 {
			destData = this.doMerge(destData, sourceData[startIndex:], i + 1 - startIndex)
		}
	}

	ret := make([]Record, len(destData))
	copy(ret, destData)
	return ret
}

func (this *periodConverter) convertSimple(sourceData []Record) []Record {
	multiplier := this.destPeriod.UnitCount() / this.srcPeriod.UnitCount()
	destData := make([]Record, 0, (len(sourceData) + multiplier - 1) / multiplier)

	destData = this.doMerge(destData, sourceData, multiplier)

	if len(destData) != cap(destData) {
		panic(errors.New("convertSimple ASSERT fail"))
	}

	return destData
}

func (this *periodConverter) Convert(sourceData []Record) []Record {
	if this.destPeriod.Eq(this.srcPeriod) {
		return sourceData
	}

	switch this.destPeriod.Unit() {
	case PERIOD_UNIT_MINUTE:
		switch this.srcPeriod.Unit() {
		case PERIOD_UNIT_MINUTE:
			return this.convertMinute2Minute(sourceData)
		default:
			panic(errors.New("bad source period"))
		}
	case PERIOD_UNIT_DAY:
		switch this.srcPeriod.Unit() {
		case PERIOD_UNIT_MINUTE:
			return this.convertMinute2Day(sourceData)
		case PERIOD_UNIT_DAY:
			return this.convertSimple(sourceData)
		default:
			panic(errors.New("bad source period"))
		}
	case PERIOD_UNIT_WEEK:
		switch this.srcPeriod.Unit() {
		case PERIOD_UNIT_DAY:
			return this.convertDay2Week(sourceData)
		case PERIOD_UNIT_WEEK:
			return this.convertSimple(sourceData)
		default:
			panic(errors.New("bad source period"))
		}
	case PERIOD_UNIT_MONTH:
		switch this.srcPeriod.Unit() {
		case PERIOD_UNIT_DAY:
			return this.convertDay2Month(sourceData)
		case PERIOD_UNIT_MONTH:
			return this.convertSimple(sourceData)
		default:
			panic(errors.New("bad source period"))
		}
	case PERIOD_UNIT_QUARTER:
		switch this.srcPeriod.Unit() {
		case PERIOD_UNIT_DAY:
			fallthrough
		case PERIOD_UNIT_MONTH:
			return this.convert2Quarter(sourceData)
		case PERIOD_UNIT_QUARTER:
			return this.convertSimple(sourceData)
		default:
			panic(errors.New("bad source period"))
		}
	case PERIOD_UNIT_YEAR:
		switch this.srcPeriod.Unit() {
		case PERIOD_UNIT_DAY:
			fallthrough
		case PERIOD_UNIT_MONTH:
			fallthrough
		case PERIOD_UNIT_QUARTER:
			return this.convert2Year(sourceData)
		case PERIOD_UNIT_YEAR:
			return this.convertSimple(sourceData)
		default:
			panic(errors.New("bad source period"))
		}
	}

	return nil
}

// Forward Adjust Price Converter


func NewForwardAdjustConverter(period Period, items []InfoExItem) Converter {
	cpy := make([]InfoExItem, len(items))
	for i, item := range items {
		cpy[i] = item
	}

	sort.SliceStable(cpy, func (i, j int) bool {
		return cpy[i].Date < cpy[j].Date
	})

	return &forwardAdjustConverter{period: period, items: cpy}
}

func (this *forwardAdjustConverter) doConvert(data []Record, item *InfoExItem) {
	var forwardAdjustPrice = func (price int32) int32 {
		fPrice := float32(price) / 1000.0
		ret := int32(((fPrice - item.Bonus) + item.RationedShares * item.RationedSharePrice) / (1 + item.DeliveredShares + item.RationedShares) * 10000.0)
		ret = (ret + 5) / 10
		return ret
	}

	for i := 0; i < len(data); i++ {
		r := &data[i]
		if GetDateDay(r.Date) >= int(item.Date) {
			break
		}

		r.Open = forwardAdjustPrice(r.Open)
		r.Close = forwardAdjustPrice(r.Close)
		r.Low = forwardAdjustPrice(r.Low)
		r.High = forwardAdjustPrice(r.High)
	}
}

func (this *forwardAdjustConverter) Convert(sourceData []Record) []Record {
	if len(sourceData) == 0 {
		return sourceData
	}

	if len(this.items) == 0 {
		return sourceData
	}

	firstDate := GetDateDay(sourceData[0].Date)
	lastDate := GetDateDay(sourceData[len(sourceData) - 1].Date)

	if int(this.items[len(this.items) - 1].Date) <= firstDate {
		return sourceData
	}

	ret := make([]Record, len(sourceData))
	copy(ret, sourceData)

	for i := range this.items {
		item := &this.items[i]
		if int(item.Date) <= firstDate {
			continue
		}
		if int(item.Date) > lastDate {
			break
		}

		this.doConvert(ret, item)
	}

	return ret
}
