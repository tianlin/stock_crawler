package stock_crawler

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestFetchInfo(t *testing.T) {
	re := regexp.MustCompile(Pattern)
	content := `var hq_str_sh510050="50ETF,0.000,2.453,2.453,0.000,0.000,0.000,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,2018-12-14,09:10:39,00";`
	if !re.MatchString(content) {
		t.Errorf("Match failed")
	}

	submatches := re.FindStringSubmatch(content)

	if submatches[2] != "510050" {
		t.Errorf("Get stock id failed:%s", submatches[1])
	}

	if submatches[3] != "50ETF,0.000,2.453,2.453,0.000,0.000,0.000,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,0,0.000,2018-12-14,09:10:39,00" {
		t.Errorf("Get info failed:%s", submatches[2])
	}
}

func TestCrawl(t *testing.T) {
	ids := `510300,510500,510050,603799,002460,002466,600362,603616,600876,000672,000786,600581,000935,000830,600160,300487,002497,600596,600036,601166,601998,000001,601628,000627,601601,601318,601336,600030`
	sc, err := NewStockCrawler(strings.Split(ids, ","))
	if err != nil {
		t.Errorf("New stock crawler failed: %s", err.Error())
	}

	infos := sc.Crawl()
	if len(infos) != len(strings.Split(ids, ",")) {
		t.Errorf("Crawl failed: %d vs %d", len(infos), len(strings.Split(ids, ",")))
	} else {
		for id, info := range infos {
			fmt.Printf("id:%s, info:%v\n", id, info)
		}
	}
}
