package fetcher

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/yagihash/fsw-calendar/fetcher/class"
	"github.com/yagihash/fsw-calendar/fetcher/course"
)

var (
	urlTmpl = "https://%s/driving/sports/%s/%s/%d/%02d.html"
)

type DocEvent struct {
	Date  string
	Start string
	End   string
	Title string
}

type Client interface {
	Get(url string) (*http.Response, error)
}

type Fetcher struct {
	hostname   string
	course     course.Course
	class      class.Class
	httpclient Client
}

func New(hostname string, course course.Course, class class.Class, c Client) *Fetcher {
	return &Fetcher{
		hostname:   hostname,
		course:     course,
		class:      class,
		httpclient: c,
	}
}

func (f *Fetcher) FetchDocEvents(y, m int) ([]DocEvent, error) {
	var rawEvents []DocEvent

	url := fmt.Sprintf(urlTmpl, f.hostname, f.course, f.class, y, m)

	res, err := f.httpclient.Get(url)
	if err != nil {
		return rawEvents, fmt.Errorf("failed to fetch raw events from %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return rawEvents, fmt.Errorf("got status code %d on %s", res.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return rawEvents, fmt.Errorf("failed to initialize document reader: %w", err)
	}

	doc.Find("#table-calendar > tbody > tr.row-rc > td.type > div > p").Each(func(i int, s *goquery.Selection) {
		d, _ := s.Parent().Parent().Parent().Attr("data-date")
		t := strings.Split(doc.Find("#table-calendar > tbody > tr.row-rc > td.time > div > p").Eq(i).Text(), "~")
		rawEvents = append(rawEvents, DocEvent{d, t[0], t[1], s.Text()})
	})

	return rawEvents, nil
}
