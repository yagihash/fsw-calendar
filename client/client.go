package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yagihash/fsw-calendar/event"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	URL string
}

func New(u string) *Client {
	return &Client{
		URL: u,
	}
}

func (c *Client) Fetch() (event.Events, error) {
	var events []*event.Event

	res, err := http.Get(c.URL)
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
		events = append(events, event.New(d, t[0], t[1], s.Text()))
	})

	return events, nil
}
