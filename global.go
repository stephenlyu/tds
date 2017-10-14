package tds

import (
	"time"
	"runtime"
)

var Local *time.Location

func init() {
	if runtime.GOOS == "windows" {
		// FIXME:
		Local = time.Local
	} else {
		Local, _ = time.LoadLocation("Asia/Shanghai")
	}
}

func SetLocationName(name string) error {
	l, err := time.LoadLocation(name)
	if err != nil {
		return err
	}
	Local = l
	return nil
}
