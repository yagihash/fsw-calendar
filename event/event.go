package event

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

const tz = "Asia/Tokyo"

type Event struct {
	*calendar.Event
}

func New(d, s, e, title string) *Event {
	template := "%sT%s:00+09:00"

	return &Event{
		&calendar.Event{
			Summary: title,
			Start: &calendar.EventDateTime{
				DateTime: fmt.Sprintf(template, d, strings.TrimPrefix(s, "0")),
			},
			End: &calendar.EventDateTime{
				DateTime: fmt.Sprintf(template, d, strings.TrimPrefix(e, "0")),
			},
		},
	}
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
