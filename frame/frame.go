package frame

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"google.golang.org/api/calendar/v3"
)

const tz = "Asia/Tokyo"

type Frame struct {
	Title string
	Start time.Time
	End   time.Time
}

func New(d, s, e, title string) *Frame {
	template := "%sT%s:00+09:00"

	start, err := time.Parse(time.RFC3339, fmt.Sprintf(template, d, s))
	if err != nil {
		return nil
	}

	end, err := time.Parse(time.RFC3339, fmt.Sprintf(template, d, e))
	if err != nil {
		return nil
	}

	f := &Frame{
		Title: title,
		Start: start,
		End:   end,
	}
	return f
}

func (f *Frame) Event() *calendar.Event {
	if f == nil {
		return nil
	}

	return &calendar.Event{
		Summary: f.Title,
		Start: &calendar.EventDateTime{
			DateTime: f.Start.Format(time.RFC3339),
			TimeZone: tz,
		},
		End: &calendar.EventDateTime{
			DateTime: f.End.Format(time.RFC3339),
			TimeZone: tz,
		},
	}
}

func Fetch(url string) ([]*Frame, error) {
	var frames []*Frame

	res, err := http.Get(url)
	if err != nil {
		return frames, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return frames, fmt.Errorf("status code error: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return frames, err
	}

	doc.Find("#table-calendar > tbody > tr.row-rc > td.type > div > p").Each(func(i int, s *goquery.Selection) {
		d, _ := s.Parent().Parent().Parent().Attr("data-date")
		t := strings.Split(doc.Find("#table-calendar > tbody > tr.row-rc > td.time > div > p").Eq(i).Text(), "~")
		frames = append(frames, New(d, t[0], t[1], s.Text()))
	})

	return frames, nil
}

func NewFromEvent(event *calendar.Event) *Frame {
	start, err := time.Parse(time.RFC3339, event.Start.DateTime)
	if err != nil {
		panic(err)
	}

	end, err := time.Parse(time.RFC3339, event.End.DateTime)
	if err != nil {
		panic(err)
	}

	return &Frame{
		Title: event.Summary,
		Start: start,
		End:   end,
	}
}
