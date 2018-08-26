package xlsx

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

func Write(xlsxFile string, data [][]interface{}) error {
	xlsx := excelize.NewFile()
	// Create a new sheet.
	sheetName := "Sheet1"
	index := xlsx.NewSheet(sheetName)

	for row, values := range data {
		for col, v := range values {
			axis := GetAxis(row, col)
			xlsx.SetCellValue(sheetName, axis, v)
		}
	}

	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	return xlsx.SaveAs(xlsxFile)
}
