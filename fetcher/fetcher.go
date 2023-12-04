package fetcher

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/yagihash/fsw-calendar/fetcher/class"
	"github.com/yagihash/fsw-calendar/fetcher/course"
	"github.com/yagihash/fsw-calendar/utils"
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

func (f *Fetcher) FetchDocEvents(y, m, length int) ([]DocEvent, error) {
	var rawEvents []DocEvent

	for i := 0; i < length; i++ {
		url := fmt.Sprintf(urlTmpl, f.hostname, f.course, f.class, y, m)

		res, err := f.httpclient.Get(url)
		if err != nil {
			return rawEvents, fmt.Errorf("failed to fetch raw events from %s: %w", url, err)
		}

		// Even if the schedule for the next month or later is not public, it is not the error.
		if i != 0 && res.StatusCode == http.StatusNotFound {
			return rawEvents, nil
		}

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

		y, m = utils.NextMonth(y, m)

		_ = res.Body.Close()
	}

	return rawEvents, nil
}
