package event

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"

	"github.com/PuerkitoBio/goquery"
)

type Events []*calendar.Event

func (es Events) Diff(another Events) (negative, positive Events) {
	for _, e := range another {
		if !es.Has(e) {
			negative = append(negative, e)
		}
	}

	for _, e := range es {
		if !another.Has(e) {
			positive = append(positive, e)
		}
	}

	return
}

func (es Events) Has(b *calendar.Event) bool {
	sb, err := time.Parse(time.RFC3339, b.Start.DateTime)
	if err != nil {
		panic(err)
	}

	eb, err := time.Parse(time.RFC3339, b.End.DateTime)
	if err != nil {
		panic(err)
	}

	for _, a := range es {
		sa, err := time.Parse(time.RFC3339, a.Start.DateTime)
		if err != nil {
			panic(err)
		}

		ea, err := time.Parse(time.RFC3339, a.End.DateTime)
		if err != nil {
			panic(err)
		}

		if a.Summary == b.Summary &&
			sa == sb &&
			ea == eb {
			return true
		}
	}

	return false
}

func (es Events) Add(e *calendar.Event) {
	es = append(es, e)
}

func Fetch(url string) (Events, error) {
	var events []*calendar.Event

	res, err := http.Get(url)
	if err != nil {
		return events, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return events, fmt.Errorf("status code error: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return events, err
	}

	doc.Find("#table-calendar > tbody > tr.row-rc > td.type > div > p").Each(func(i int, s *goquery.Selection) {
		d, _ := s.Parent().Parent().Parent().Attr("data-date")
		t := strings.Split(doc.Find("#table-calendar > tbody > tr.row-rc > td.time > div > p").Eq(i).Text(), "~")
		events = append(events, New(d, t[0], t[1], s.Text()))
	})

	return events, nil
}
