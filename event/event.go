package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/yagihash/fsw-calendar/fetcher"

	"google.golang.org/api/calendar/v3"
)

type Event struct {
	*calendar.Event
}

func New(date, start, end, title string) *Event {
	template := "%sT%s:00+09:00"

	return &Event{
		&calendar.Event{
			Summary: title,
			Start: &calendar.EventDateTime{
				DateTime: fmt.Sprintf(template, date, strings.TrimPrefix(start, "0")),
			},
			End: &calendar.EventDateTime{
				DateTime: fmt.Sprintf(template, date, strings.TrimPrefix(end, "0")),
			},
		},
	}
}

func NewFromDocEvent(d fetcher.DocEvent) *Event {
	return New(d.Date, d.Start, d.End, d.Title)
}

func (e *Event) Equals(c *Event) bool {
	es, err := time.Parse(time.RFC3339, e.Start.DateTime)
	if err != nil {
		panic(err)
	}

	ee, err := time.Parse(time.RFC3339, e.End.DateTime)
	if err != nil {
		panic(err)
	}

	cs, err := time.Parse(time.RFC3339, c.Start.DateTime)
	if err != nil {
		panic(err)
	}

	ce, err := time.Parse(time.RFC3339, c.End.DateTime)
	if err != nil {
		panic(err)
	}

	return e.Summary == c.Summary && es.Equal(cs) && ee.Equal(ce)
}
