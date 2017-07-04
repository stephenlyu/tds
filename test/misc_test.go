package test

import (
. "github.com/onsi/ginkgo"

	"fmt"
	"regexp"
	"time"
	"encoding/gob"
	"bytes"
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

var _ = Describe("TestSlice", func() {
	It("test", func (){
		var s = make([]int, 0, 100)
		fmt.Printf("%d, %d, %p\n", len(s), cap(s), s)
		for i := 0; i < 100; i++ {
			s = append(s, i)
			fmt.Printf("%d, %d, %p\n", len(s), cap(s), s)
		}
	})
})

var _ = Describe("TestDate", func() {
	It("test", func (){
		now := time.Now()
		now.Weekday()
	})
})

var _ = Describe("TestGob", func() {
	It("test", func (){
		type A struct {
			I int
			F float32
		}

		var m = map[string][]A {
			"600000": []A{
				A{I: 100, F: 1.34},
				A{I: 200, F: 3.89},
			},
		}
		buffer := &bytes.Buffer{}
		encoder := gob.NewEncoder(buffer)
		encoder.Encode(m)
		fmt.Println(buffer.Bytes())

		decoder := gob.NewDecoder(bytes.NewBuffer(buffer.Bytes()))
		var m1 map[string][]A
		err := decoder.Decode(&m1)
		fmt.Println(err, m1)
	})
})

var _ = Describe("TestConv", func() {
	It("test", func (){
		i := -1
		ui := uint32(i)
		fmt.Println(ui)
	})
})
