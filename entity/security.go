package entity

import (
	"regexp"
	"errors"
	"fmt"
)

type Security struct {
	Code string
	Exchange string
}

var codePattern, _ = regexp.Compile(`^([0-9a-zA-Z]+)\.([0-9a-zA-Z]+)$`)

func ParseSecurity(securityCode string) (*Security, error) {
	matches := codePattern.FindSubmatch([]byte(securityCode))
	if len(matches) == 0 {
		return nil, errors.New("Bad security code")
	}

	code := string(matches[1])
	exchange := string(matches[2])
	return &Security{code, exchange}, nil
}

func (this *Security) String() string {
	return fmt.Sprintf("%s.%s", this.Code, this.Exchange)
}
