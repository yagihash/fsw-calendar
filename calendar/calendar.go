package calendar

import (
	"context"
	"time"

	"github.com/yagihash/fsw-calendar/event"
	"github.com/yagihash/fsw-calendar/utils"

	"google.golang.org/api/calendar/v3"
)

type Calendar struct {
	cs         *calendar.Service
	calendarID string
	tz         *time.Location
}

func New(ctx context.Context, calendarID string, tz *time.Location) (*Calendar, error) {
	cs, err := calendar.NewService(ctx)
	if err != nil {
		return nil, err
	}

	cal := &Calendar{
		cs:         cs,
		calendarID: calendarID,
		tz:         tz,
	}

	return cal, nil
}

func (cal *Calendar) GetEvents(y, m, length int) (event.Events, error) {
	start := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, cal.tz).Format(time.RFC3339)
	end := start

	for i, tmpY, tmpM := 0, y, m; i < length; i++ {
		tmpY, tmpM = utils.NextMonth(tmpY, tmpM)
		end = time.Date(tmpY, time.Month(tmpM), 1, 0, 0, 0, 0, cal.tz).Format(time.RFC3339)
	}

	events, err := cal.cs.Events.List(cal.calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(start).TimeMax(end).Do()
	if err != nil {
		return event.Events{}, err
	}

	e := event.NewEvents(events.Items)

	return e, nil
}

func (cal *Calendar) Insert(e *calendar.Event) error {
	_, err := cal.cs.Events.Insert(cal.calendarID, e).Do()
	return err
}

func (cal *Calendar) Delete(EventID string) error {
	return cal.cs.Events.Delete(cal.calendarID, EventID).Do()
}
