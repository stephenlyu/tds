package tds


type Converter interface {
	Convert(sourceData []Record) []Record
}

type periodConverter struct {
	srcPeriod Period
	destPeriod Period
}

type forwardAdjustConverter struct {
	items []InfoExItem
}

// Period Data Converters

func NewPeriodConverter(srcPeriod Period, destPeriod Period) Converter {
	return &periodConverter{srcPeriod: srcPeriod, destPeriod: destPeriod}
}

func (this *periodConverter) Convert(sourceData []Record) []Record {
	return nil
}

// Forward Adjust Price Converter


func NewForwardAdjustConverter(items []InfoExItem) Converter {
	return &forwardAdjustConverter{items: items}
}

func (this *forwardAdjustConverter) Convert(sourceData []Record) []Record {
	return nil
}
