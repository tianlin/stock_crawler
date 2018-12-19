package stock_crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	BatchSize = 10
	Pattern   = `var hq_str_(?P<StockID>[^=]+)="(?P<Info>[^"]*)";`
)

type StockInfo struct {
	Price float64
	Date  string
	Time  string
	Id    string
}

type StockCrawler struct {
	ids    []string
	re     *regexp.Regexp
	client *http.Client
}

func NewStockCrawler(ids []string) (sc *StockCrawler, err error) {
	if len(ids) == 0 {
		err = fmt.Errorf("ids is empty")
		return
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	sc = &StockCrawler{
		ids:    ids,
		re:     regexp.MustCompile(Pattern),
		client: client,
	}
	return
}

func updateMap(stocks map[string]*StockInfo, inChan chan *StockInfo, stopChan chan struct{}) {
	for {
		select {
		case info := <-inChan:
			stocks[info.Id] = info
		case <-stopChan:
			close(stopChan)
			return
		}
	}
}

func (sc *StockCrawler) Crawl() (stocks map[string]*StockInfo) {
	stocks = map[string]*StockInfo{}
	stopChan := make(chan struct{})
	infoChan := make(chan *StockInfo)
	go updateMap(stocks, infoChan, stopChan)

	wg := sync.WaitGroup{}
	for i := 0; i < len(sc.ids)/BatchSize+1; i++ {
		wg.Add(1)
		begin := i * BatchSize
		end := (i + 1) * BatchSize
		if len(sc.ids) < end {
			end = len(sc.ids)
		}
		go func(ids []string) {
			defer wg.Done()
			var idsWithPrefix []string
			for _, id := range ids {
				prefix := []rune(id)[0]
				if unicode.IsDigit(prefix) {
					if prefix == '0' || prefix == '3' {
						idsWithPrefix = append(idsWithPrefix, "sz"+id)
					} else {
						idsWithPrefix = append(idsWithPrefix, "sh"+id)
					}
				} else {
					idsWithPrefix = append(idsWithPrefix, id)
				}
			}
			url := fmt.Sprintf("http://hq.sinajs.cn/list=%s", strings.Join(idsWithPrefix, ","))
			resp, err := sc.client.Get(url)
			if err != nil {
				log.Printf("Get %s failed: %s\n", url, err.Error())
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Printf("Get %s failed: %s\n", url, resp.Status)
				return
			}

			content, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Read %s failed: %s", url, err.Error())
			}

			for _, line := range strings.Split(string(content), "\n") {
				if line == "" {
					continue
				}

				matches := sc.re.FindStringSubmatch(line)
				if len(matches) > 2 {
					stockId := matches[1]
					stockId = strings.TrimLeft(stockId, "sh")
					stockId = strings.TrimLeft(stockId, "sz")
					if matches[2] == "" {
						infoChan <- &StockInfo{
							Id: stockId,
						}
						continue
					}

					infos := strings.Split(matches[2], ",")
					if len(infos) > 6 {
						price, _ := strconv.ParseFloat(infos[3], 10)
						info := StockInfo{
							Id:    stockId,
							Price: price,
							Date:  infos[len(infos)-3],
							Time:  infos[len(infos)-2],
						}

						infoChan <- &info
					} else {
						log.Printf("%s Line is invalid: %s", url, line)
					}
				} else {
					log.Printf("%s Line is invalid: %s", url, line)
				}
			}
		}(sc.ids[begin:end])
	}

	wg.Wait()
	stopChan <- struct{}{}
	<-stopChan

	return
}
