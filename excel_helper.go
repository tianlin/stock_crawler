package stock_crawler

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"time"
)

const (
	DateFmt   = "2006/1/2"
	TimeFmt   = "1504"
	SheetName = "Sheet1"
)

func GetIds(filename string) (ids []string, err error) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		err = fmt.Errorf("Open %s failed: %s", filename, err.Error())
		return
	}

	rows, err := xlsx.Rows(SheetName)
	i := 0
	for ; rows.Next(); i++ {
		if i == 3 {
			cols := rows.Columns()
			if len(cols) < 5 {
				err = fmt.Errorf("%s is invalid: columns is too less", filename)
				return
			}

			for _, colCell := range cols[4:] {
				if colCell != "" {
					ids = append(ids, colCell)
				}
			}
			break
		} else {
			continue
		}

	}

	if i != 3 {
		err = fmt.Errorf("%s is invalid: rows is too less", filename)
		return
	}

	return
}

func UpdateInfos(filename string, ids []string, stocks map[string]*StockInfo, now time.Time) (err error) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		err = fmt.Errorf("Open %s failed: %s", filename, err.Error())
		return
	}

	// TODO: if row num can be calculate by now, remove this to be faster.
	rows, err := xlsx.Rows(SheetName)
	j := 0
	for ; rows.Next(); j++ {
		cols := rows.Columns()
		if j >= 4 && cols[0] == "" {
			break
		}
	}

	if j < 4 {
		err = fmt.Errorf("%s is invalid: rows is too less", filename)
		return
	}

	// Add Date & Time
	xlsx.SetCellStr(SheetName, fmt.Sprintf("A%d", j+1), now.Format(DateFmt))
	xlsx.SetCellStr(SheetName, fmt.Sprintf("B%d", j+1), now.Format(TimeFmt))

	// Add current prices
	for i, id := range ids {
		alpha := excelize.ToAlphaString(i + 4)
		axis := fmt.Sprintf("%s%d", alpha, j+1)
		if info, ok := stocks[id]; ok && info.Price != 0 {
			xlsx.SetCellValue(SheetName, axis, info.Price)
		} else {
			// Use last price as current price
			price := float64(0)
			if info.Price == 0 && j > 4 {
				lastAxis := fmt.Sprintf("%s%d", alpha, j)
				lastPrice := xlsx.GetCellValue(SheetName, lastAxis)
				price, _ = strconv.ParseFloat(lastPrice, 64)
			}
			xlsx.SetCellValue(SheetName, axis, price)
		}
	}

	return xlsx.Save()
}
