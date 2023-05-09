package function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"google.golang.org/api/calendar/v3"

	"github.com/yagihash/fsw-calendar/config"
	"github.com/yagihash/fsw-calendar/event"
	"github.com/yagihash/fsw-calendar/logger"
)

func Register(ctx context.Context, message *pubsub.Message) error {
	log, err := logger.New()
	if err != nil {
		return err
	}

	defer func() { _ = log.Sync() }()

	log.Debug("logger is ready")

	c, err := config.Load()
	if err != nil {
		log.Error("failed to load config", zap.Error(err))
		return err
	}

	var data config.Data
	if err := json.Unmarshal(message.Data, &data); err != nil {
		log.Error("failed to unmarshal message", zap.Error(err), zap.Any("data", message.Data))
	}

	jst, err := time.LoadLocation(c.Timezone)
	if err != nil {
		log.Error("failed to load timezone", zap.Error(err))
		return err
	}

	y := time.Now().In(jst).Year()
	m := int(time.Now().In(jst).Month())

	fetchedEvents, err := FetchEvents(data.URL, y, m, c.Recurrence)
	if err != nil {
		log.Error(
			"failed to fetch schedule data",
			zap.Error(err),
			zap.String("url", data.URL),
			zap.Int("y", y),
			zap.Int("m", m),
			zap.Int("recurrence", c.Recurrence),
		)
	}

	log.Info("loaded schedules", zap.Any("events", fetchedEvents))

	cs, err := calendar.NewService(ctx)
	if err != nil {
		log.Error("failed to access google calendar API", zap.Error(err))
		return err
	}

	existingEvents, err := ListExistingEvents(cs, y, m, c.Recurrence, jst, data.CalendarID)
	if err != nil {
		log.Error("failed to load existing events from google calendar", zap.Error(err))
		return err
	}

	toBeAdded, toBeDeleted := existingEvents.Diff(fetchedEvents)
	if len(toBeAdded) == 0 && len(toBeDeleted) == 0 {
		log.Debug("no update", zap.Any("existing", existingEvents), zap.Any("fetched", fetchedEvents))
	}

	for _, e := range toBeAdded {
		if e == nil {
			continue
		}

		_, err := cs.Events.Insert(data.CalendarID, e.Event).Do()
		if err != nil {
			log.Error("failed to insert event", zap.Error(err), zap.Any("event", e))
		} else {
			log.Info("added new event", zap.Any("event", e))
		}
	}

	for _, e := range toBeDeleted {
		if err := cs.Events.Delete(data.CalendarID, e.Id).Do(); err != nil {
			log.Error("failed to reset event", zap.Any("event", e))
		} else {
			log.Info("deleted stale event", zap.Any("event", e))
		}
	}

	return nil
}

func NextMonth(y, m int) (int, int) {
	var nextY, nextM int

	if m == 12 {
		nextY = y + 1
		nextM = 1
	} else {
		nextY = y
		nextM = m + 1
	}

	return nextY, nextM
}

func FetchEvents(tmpl string, y, m, rec int) (event.Events, error) {
	events := event.Events{}
	for i := 0; i < rec; i++ {
		url := fmt.Sprintf(tmpl, y, m)

		tmp, err := event.Fetch(url)
		if err != nil {
			return events, err
		}

		events = append(events, tmp...)
		y, m = NextMonth(y, m)
	}

	return events.Unique(), nil
}

func ListExistingEvents(cs *calendar.Service, y, m, rec int, tz *time.Location, calendarID string) (event.Events, error) {
	start := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, tz).Format(time.RFC3339)
	end := start

	for i, tmpY, tmpM := 0, y, m; i < rec; i++ {
		tmpY, tmpM = NextMonth(tmpY, tmpM)
		fmt.Println(tmpY, tmpM)
		end = time.Date(tmpY, time.Month(tmpM), 1, 0, 0, 0, 0, tz).Format(time.RFC3339)
	}

	events, err := cs.Events.List(calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(start).TimeMax(end).Do()
	if err != nil {
		return event.Events{}, err
	}

	e := event.NewEvents(events.Items)

	return e.Unique(), nil
}
