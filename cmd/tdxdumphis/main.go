package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/datasource/tdx"
)

func chk(err error) {
	if err != nil {
		fmt.Printf("[ERROR] error: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	periodStr := flag.String("period", "D1", "Period to get")
	dataDir := flag.String("data-dir", "data", "Data directory")
	flag.Parse()

	err, dp := period.PeriodFromString(*periodStr)
	chk(err)

	if len(flag.Args()) == 0 {
		return
	}

	code := flag.Args()[0]
	security, err := entity.ParseSecurity(code)
	chk(err)

	ds := tdxdatasource.NewDataSource(*dataDir, true)
	err, records := ds.GetData(security, dp)
	chk(err)

	for i := range records {
		fmt.Printf("%+v\n", &records[i])
	}
}
