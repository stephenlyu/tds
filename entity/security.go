package entity

import (
	"regexp"
	"errors"
	"strings"
	"github.com/stephenlyu/tds/util"
	"fmt"
)

type Security struct {
	Category string
	Code     string
	Exchange string

	fullCode string
}

var (
	aShareCodePattern, _ = regexp.Compile(`^()([0-9]+)\.([0-9a-zA-Z]+)$`)
	okexFutureCodePattern, _ = regexp.Compile(`^([A-Z]+)([TNQ]FUT)\.(OKEX)$`)
	okexSwapCodePattern, _ = regexp.Compile(`^([A-Z]+)(FUT[F]?)\.(OKEX)$`)
	ploFutureCodePattern, _ = regexp.Compile(`^([A-Z]+)(FUT|INDEX)\.(PLO)$`)
	bitmexSwapCodePattern, _ = regexp.Compile(`^([A-Z]+)(FUT)\.(BITMEX)$`)
	bitmexFutureCodePattern, _ = regexp.Compile(`^([A-Z]+)([A-Z][0-9]+)\.(BITMEX)$`)
	dcSpotCodePattern, _  = regexp.Compile(`^([A-Z]+_[A-Z]+)(SPOT)\.([A-Z]+)$`)
	cnCommodityFutureCodePattern, _ = regexp.Compile(`^([A-Z]+)([0-9]+)\.([A-Z]+)$`)
	indexPattern, _ = regexp.Compile(`^([A-Z]+)(FUT|INDEX)\.([A-Z]+)$`)
	codePatterns = []*regexp.Regexp{aShareCodePattern, okexFutureCodePattern,
		okexSwapCodePattern, bitmexSwapCodePattern, bitmexFutureCodePattern,
		ploFutureCodePattern,
		dcSpotCodePattern,
		cnCommodityFutureCodePattern,
		indexPattern,
	}
)

func ParseSecurity(securityCode string) (*Security, error) {
	securityCode = strings.ToUpper(securityCode)

	var matches [][]byte
	for _, p := range codePatterns {
		matches = p.FindSubmatch([]byte(securityCode))
		if len(matches) > 0 {
			break
		}
	}

	if len(matches) == 0 {
		return nil, errors.New("Bad security code")
	}

	var cat, code, exchange string
	switch len(matches) {
	case 3:
		code = string(matches[1])
		exchange = string(matches[2])
	case 4:
		cat = string(matches[1])
		code = string(matches[2])
		exchange = string(matches[3])
	}

	return &Security{cat, code, exchange, securityCode}, nil
}

func ParseSecurityUnsafe(securityCode string) (*Security) {
	ret, _ := ParseSecurity(securityCode)
	return ret
}

func (this *Security) String() string {
	return this.fullCode
}

func (this *Security) CatCode() string {
	return fmt.Sprintf("%s.%s", this.GetCategory(), this.Exchange)
}

func (this *Security) GetCategory() string {
	if this.Category == "" {
		return "ASTOCK"
	}
	return this.Category
}

func (this *Security) GetCode() string {
	return this.Code
}

func (this *Security) GetExchange() string {
	return this.Exchange
}

func (this *Security) IsSpot() bool {
	return this.Code == "SPOT"
}

func (this *Security) IsIndex() bool {
	return this.Code == "INDEX"
}

func (this *Security) IsDigitCurrency() bool {
	return util.InStrings(this.Exchange, []string{"OKEX", "HUOBI"})
}
