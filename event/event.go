package event

import (
	"fmt"
	"strings"

	"google.golang.org/api/calendar/v3"
)

const tz = "Asia/Tokyo"

func New(d, s, e, title string) *calendar.Event {
	template := "%sT%s:00+09:00"

	return &calendar.Event{
		Summary: title,
		Start: &calendar.EventDateTime{
			DateTime: fmt.Sprintf(template, d, strings.TrimPrefix(s, "0")),
		},
		End: &calendar.EventDateTime{
			DateTime: fmt.Sprintf(template, d, strings.TrimPrefix(e, "0")),
		},
	}
}
