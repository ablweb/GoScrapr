package wikixls

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/ablweb/GoScrapr/pkg/scraper"
	"github.com/gocolly/colly/v2"
	"github.com/xuri/excelize/v2"
)

var (
	titleStyleDef = excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  14,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#000000"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 2},
			{Type: "right", Color: "#000000", Style: 2},
			{Type: "bottom", Color: "#000000", Style: 2},
		},
	}

	headerStyleDef = excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#D9E1F2"}, // Light blue header
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	}

	cellStyleDef = excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	}
)

func CreateTableFrom(file *excelize.File, url string, tableSelector string) error {
	reach, err := scraper.IsReachable(url)
	if !reach {
		return err
	}

	titleStyle, _ := file.NewStyle(&titleStyleDef)
	headerStyle, _ := file.NewStyle(&headerStyleDef)
	cellStyle, _ := file.NewStyle(&cellStyleDef)

	c := colly.NewCollector()

	var tableOffset int
	var tableId int = 1

	c.OnHTML(tableSelector, func(table *colly.HTMLElement) {
		var tLength int
		var maxCol int
		table.DOM.Find("tr").Each(func(row int, tr *goquery.Selection) {
			tLength++
			var currentMax int
			tr.Find("th,td").Each(func(col int, td *goquery.Selection) {
				tag := goquery.NodeName(td)
				val := td.Text()

				cell := toExcelCell(col, row+tableOffset+1)
				file.SetCellValue("Sheet1", cell, val)
				if tag == "td" {
					file.SetCellStyle("Sheet1", cell, cell, cellStyle)
				} else {
					file.SetCellStyle("Sheet1", cell, cell, headerStyle)
				}
				currentMax++
			})
			if currentMax > maxCol {
				maxCol = currentMax
			}
		})

		title := fmt.Sprintf("Table %d — %s", tableId, url)
		titleCell := toExcelCell(maxCol/2, tableOffset)
		startCell := toExcelCell(maxCol-1, tableOffset)
		endCell := toExcelCell(0, tableOffset)
		file.SetCellValue("Sheet1", titleCell, title)
		file.SetCellStyle("Sheet1", startCell, endCell, titleStyle)

		tableOffset += tLength + 2
		tableId++
	})

	// error while scraping page
	if err := c.Visit(url); err != nil {
		return err
	}
	return nil
}

func toExcelCell(col, row int) string {
	colName := ""
	for col >= 0 {
		colName = string('A'+(col%26)) + colName
		col = col/26 - 1
	}
	return colName + strconv.Itoa(row+1) // Excel rows start at 1
}
