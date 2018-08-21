package storage

import (
	"testing"
	"reflect"
	"fmt"
)

type A struct {
	A int
	B float64
	C string
	D float32
	E uint16
	F int32
}

func TestCsvEngine_Save(t *testing.T) {
	records := []interface{} {
		&A{1, 3.0, "first", 9.8, 255, 888},
		&A{2, 8.0, "second\"", 19.8, 2555, 999},
	}
	e := NewCsvEngine(reflect.TypeOf(A{}))
	err := e.Save("test.csv", records)
	fmt.Println(err)

}

func TestCsvEngine_Load(t *testing.T) {
	e := NewCsvEngine(reflect.TypeOf(A{}))
	err, data := e.Load("test.csv")
	fmt.Println(err, data)
}