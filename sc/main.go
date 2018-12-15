package main

import (
	"flag"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tianlin/stock_crawler"
	"log"
	"os"
	"time"
)

func main() {
	input := flag.String("input", "", "xlsx with stock ids")
	output := flag.String("output", "", "xlsx to store stock prices")
	flag.Parse()

	now := time.Now()
	log.Printf("Begin get stock ids from input")
	ids, err := stock_crawler.GetIds(*input)
	if err != nil {
		log.Fatalf("Get ids from input %s failed: %s", *input, err.Error())
	}

	log.Printf("Begin crawl stock prices")
	sc, err := stock_crawler.NewStockCrawler(ids)
	if err != nil {
		log.Fatalf("New stock crawler failed: %s", err.Error())
	}
	stacks := sc.Crawl()

	f, err := os.Open(*output)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Begin build output from input xlsx")
			// Copy input as output
			xlsx, _ := excelize.OpenFile(*input)
			xlsx.SaveAs(*output)
		} else {
			log.Fatalf("Open %s failed: %s", *output, err.Error())
		}
	}
	f.Close()

	log.Printf("Begin update stock prices to output")
	err = stock_crawler.UpdateInfos(*output, ids, stacks, now)
	if err != nil {
		log.Fatalf("Update infos failed: %s", err.Error())
	}

	log.Printf("Update stack prices success")
}
