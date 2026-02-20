package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ExportExcel(results []ScanResult) error {

	f := excelize.NewFile()
	sheet := "Scan Results"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"IP", "Port", "Service", "Risk"}

	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, r := range results {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.IP)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.Port)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.Service)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.Risk)
	}

	return f.SaveAs("scan_results.xlsx")
}
