package event

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Events []*Event

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

func (es Events) Has(b *Event) bool {
	for _, a := range es {
		if a.Equals(b) {
			return true
		}
	}

	return false
}

func (es Events) Unique() (unique Events) {
	for i, e := range es {
		// note: used in the case that the original calendar is broken. no need to ensure uniqueness seriously.
		if es[i+1:].Has(e) {
			// do nothing
		} else {
			unique = append(unique, e)
		}
	}

	return
}

func Fetch(url string) (Events, error) {
	var events []*Event

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
