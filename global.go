package tds

import "time"

var Local, _ = time.LoadLocation("Asia/Shanghai")

func SetLocationName(name string) error {
	l, err := time.LoadLocation(name)
	if err != nil {
		return err
	}
	Local = l
	return nil
}
