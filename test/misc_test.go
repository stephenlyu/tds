package test

import (
. "github.com/onsi/ginkgo"

	"fmt"
	"regexp"
)

var _ = Describe("TestRegexp", func() {
	It("test", func (){
		regExp, err := regexp.Compile("^([A-Z]+)([0-9]+)$")
		if err != nil {
			panic(err)
		}

		result := regExp.FindSubmatch([]byte("MINUTE10"))
		if result != nil {
			for _, r := range result {
				fmt.Println(string(r))
			}
		}
	})
})

var _ = Describe("TestParseMinute", func() {
	It("test", func (){
		regExp, err := regexp.Compile("^([0-9]{4})([0-9]{2})([0-9]{2}) ([0-9]{2}):([0-9]{2}):([0-9]{2})$")
		if err != nil {
			panic(err)
		}

		result := regExp.FindSubmatch([]byte("20170701 09:30:00"))
		if result != nil {
			for _, r := range result {
				fmt.Println(string(r))
			}
		}
	})
})
