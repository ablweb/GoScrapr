package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ablweb/GoScrapr/pkg/wikixls"
	"github.com/xuri/excelize/v2"
)

const HELP = "Usage: wikixls [URL]\n"

func run(out, errOut io.Writer, args []string) int {
	// No URL
	if len(args) < 1 {
		fmt.Fprintln(errOut, "Error: No URL to scrap")
		fmt.Fprint(out, HELP)
		return 1
	}
	url := args[0]

	file := excelize.NewFile()

	err := wikixls.CreateTableFrom(file, url, ".wikitable")
	if err != nil {
		fmt.Fprintln(errOut, "Error: While scrapping", err)
		return 1
	}

	// Save the file
	if err := file.SaveAs("example.xlsx"); err != nil {
		fmt.Println(errOut, err)
		return 0
	}

	fmt.Println("File saved successfully!")
	return 0
}

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args[1:]))
}
