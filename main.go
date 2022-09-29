package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	ExitOK = iota
	ExitError
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	u := "https://www.fsw.tv/driving/sports/ss/ss-4/2022/09.html"

	res, err := http.Get(u)
	if err != nil {
		log.Println(err)
		return ExitError
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return ExitError
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	doc.Find("#table-calendar > tbody > tr.row-rc > td.type > div > p").Each(func(i int, s *goquery.Selection) {
		d, _ := s.Parent().Parent().Parent().Attr("data-date")
		t := strings.Split(doc.Find("#table-calendar > tbody > tr.row-rc > td.time > div > p").Eq(i).Text(), "~")
		fmt.Printf("%s, %s, %s, %s\n", d, s.Text(), t[0], t[1])
	})

	return ExitOK
}
