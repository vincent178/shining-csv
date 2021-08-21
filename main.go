package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	sheetname = "Sheet1"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf(`Usage: 
  shining-csv <path to xlsx file>
`)
		os.Exit(1)

	}

	xlsxfile := os.Args[1]

	log.Println("process", xlsxfile)

	f, err := excelize.OpenFile(xlsxfile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("read", sheetname)

	// Shining always use Sheet1
	rows, err := f.GetRows(sheetname)
	if err != nil {
		log.Fatal(err)
	}

	today := time.Now().Format("2006-01-02")
	outname := today + ".csv"

	out, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}

	cw := csv.NewWriter(out)

	buf := make([]string, 0)

	buf = append(buf, "code")

	for _, row := range rows {
		for idx, col := range row {
			if idx == 0 {
				continue
			}
			buf = append(buf, col)
		}
	}

	cw.Write(buf)

	cw.Flush()
	log.Println("generated", outname)
	log.Printf("total %d lines\n", len(buf))
	log.Println("please confirm with Shining")
}
