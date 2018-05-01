package entity

import (
	"regexp"
	"errors"
	"strings"
)

type Security struct {
	Code string
	Exchange string

	fullCode string
}

var codePattern, _ = regexp.Compile(`^([0-9a-zA-Z]+)\.([0-9a-zA-Z]+)$`)

func ParseSecurity(securityCode string) (*Security, error) {
	matches := codePattern.FindSubmatch([]byte(securityCode))
	if len(matches) == 0 {
		return nil, errors.New("Bad security code")
	}

	code := string(matches[1])
	exchange := strings.ToUpper(string(matches[2]))

	return &Security{code, exchange, securityCode}, nil
}

func ParseSecurityUnsafe(securityCode string) (*Security) {
	ret, _ := ParseSecurity(securityCode)
	return ret
}

func (this *Security) String() string {
	return this.fullCode
}

func (this *Security) GetCode() string {
	return this.Code
}

func (this *Security) GetExchange() string {
	return this.Exchange
}
