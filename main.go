package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	sheetname = "Sheet1"
	batchSize = 6000

	version = "0.0.2"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf(`shinning-csv version: ` + version + `
Usage: 
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
	batch := 0
	lines := 1 // the line with header `code`

	outname := fmt.Sprintf("%s-%d.csv", today, batch)

	out, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}

	cw := csv.NewWriter(out)

	cw.Write([]string{"code"})

	includeFirstColumn := false

	for idx, row := range rows {
		for idy, col := range row {

			// 如果首行首列的值不为空，那么保存第一列的值
			if idx == 0 && idy == 0 {
				if col != "" {
					includeFirstColumn = true
				}
			}

			// 如果是不用包含首列的值则跳过
			if idy == 0 && !includeFirstColumn {
				continue
			}

			col = strings.TrimSpace(col)

			if col == "" {
				continue
			}

			lines += 1
			cw.Write([]string{col})

			if lines >= batchSize {
				// flush current file to disk
				cw.Flush()
				log.Println("generated", outname)

				lines = 1
				batch += 1

				outname = fmt.Sprintf("%s-%d.csv", today, batch)

				out, err := os.Create(outname)
				if err != nil {
					log.Fatal(err)
				}

				cw = csv.NewWriter(out)

				cw.Write([]string{"code"})
			}
		}
	}

	// flush current file to disk
	cw.Flush()
	log.Println("generated", outname)

	log.Printf("total %d lines\n", lines+batch*batchSize)
	log.Println("please confirm with Shining")
}
