package frame

import (
	"fmt"
	"time"

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
