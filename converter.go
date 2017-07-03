package tds

import "sort"

type Converter interface {
	Convert(sourceData []Record) []Record
}

type periodConverter struct {
	srcPeriod Period
	destPeriod Period
}

type forwardAdjustConverter struct {
	period Period
	items []*InfoExItem
}

// Period Data Converters

func NewPeriodConverter(srcPeriod Period, destPeriod Period) Converter {
	return &periodConverter{srcPeriod: srcPeriod, destPeriod: destPeriod}
}

func (this *periodConverter) Convert(sourceData []Record) []Record {
	return nil
}

// Forward Adjust Price Converter


func NewForwardAdjustConverter(period Period, items []InfoExItem) Converter {
	cpy := make([]*InfoExItem, len(items))
	for i, item := range items {
		cpy[i] = &item
	}

	sort.SliceStable(cpy, func (i, j int) bool {
		return cpy[i].Date < cpy[j].Date
	})

	return &forwardAdjustConverter{period: period, items: cpy}
}

func (this *forwardAdjustConverter) doConvert(data []Record, item *InfoExItem) {
	var forwardAdjustPrice = func (price int32) int32 {
		fPrice := float32(price) / 1000.0
		return int32(((fPrice - item.Bonus) + item.RationedShares * item.RationedSharePrice) / (1 + item.DeliveredShares + item.DeliveredShares) * 1000.0)
	}

	for i := 0; i < len(data); i++ {
		r := &data[i]

		if GetDateDay(this.period, r.Date) >= item.Date {
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

	firstDate := GetDateDay(this.period, sourceData[0].Date)
	lastDate := GetDateDay(this.period, sourceData[len(sourceData) - 1].Date)

	if this.items[len(this.items) - 1].Date <= firstDate {
		return sourceData
	}

	ret := make([]Record, len(sourceData))
	copy(ret, sourceData)

	for _, item := range this.items {
		if item.Date <= firstDate {
			continue
		}
		if item.Date > lastDate {
			break
		}

		this.doConvert(ret, item)
	}

	return ret
}
