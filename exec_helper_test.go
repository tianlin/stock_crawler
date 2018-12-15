package stock_crawler

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetIds(t *testing.T) {
	ids, err := GetIds("testdata/test.xlsx")
	if err != nil {
		t.Errorf("GetIds failed: %s", err.Error())
	}

	if len(ids) == 0 {
		t.Errorf("len(ids) is %d", len(ids))
	} else {
		fmt.Println(strings.Join(ids, ","))
	}
}

func TestWriteExcel(t *testing.T) {
	ids, err := GetIds("testdata/out.xlsx")
	if err != nil {
		t.Errorf("Get ids failed: %s", err.Error())
	}
	sc, err := NewStockCrawler(ids)
	if err != nil {

	}
	stacks := sc.Crawl()

	tmpFile, _ := os.Create("testdata/out.xlsx.tmp")
	file, _ := os.Open("testdata/out.xlsx")
	io.Copy(tmpFile, file)
	err = UpdateInfos("testdata/out.xlsx.tmp", ids, stacks, time.Now())
	if err != nil {
		t.Errorf("Update infos failed: %s", err.Error())
	}
	xlsx, _ := excelize.OpenFile("testdata/out.xlsx.tmp")
	value := xlsx.GetCellValue(SheetName, "A5")
	if value == "" {
		t.Errorf("Update info failed")
	}
}
