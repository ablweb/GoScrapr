package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Create a new Excel file
	f := excelize.NewFile()

	// Set values in a cell
	f.SetCellValue("Sheet1", "A1", "Hello")
	f.SetCellValue("Sheet1", "B1", "World")

	// Save the file
	if err := f.SaveAs("example.xlsx"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("File saved successfully!")
}
