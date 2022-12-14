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

	defer log.Sync()

	log.Info("logger is ready")

	c, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
		return err
	}

	var data config.Data
	json.Unmarshal(message.Data, &data)

	jst, err := time.LoadLocation(c.Timezone)
	if err != nil {
		log.Fatal("failed to load timezone", zap.Error(err))
		return err
	}

	y := time.Now().In(jst).Year()
	m := int(time.Now().In(jst).Month())

	for i := 0; i < c.Recurrence; i++ {
		url := fmt.Sprintf(data.URL, y, m)

		fetchedEvents, err := event.Fetch(url)
		if err != nil {
			log.Fatal("failed to fetch schedule data", zap.Error(err), zap.String("url", url))
			return err
		}

		log.Info("loaded schedules", zap.String("url", url))

		cs, err := calendar.NewService(ctx)
		if err != nil {
			log.Fatal("failed to access google calendar API", zap.Error(err))
			return err
		}

		nextY, nextM := NextMonth(y, m)

		events, err := cs.Events.List(data.CalendarID).ShowDeleted(false).SingleEvents(true).
			TimeMin(time.Date(y, time.Month(m), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).
			TimeMax(time.Date(nextY, time.Month(nextM), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).Do()
		if err != nil {
			log.Error("err", zap.Error(err))
			return err
		}

		existingEvents := event.Events(events.Items)

		toBeAdded, toBeDeleted := existingEvents.Diff(fetchedEvents)
		if toBeAdded == nil && toBeDeleted == nil {
			log.Info("no update", zap.Int("year", y), zap.Int("month", m))
			continue
		}

		if toBeAdded != nil {
			for _, e := range toBeAdded {
				if e == nil {
					continue
				}

				_, err := cs.Events.Insert(data.CalendarID, e).Do()
				if err != nil {
					log.Error("failed to insert event", zap.Error(err), zap.Any("event", e), zap.Int("year", y), zap.Int("month", m))
					return err
				}
			}
			log.Info("added new events", zap.Int("count", len(toBeAdded)))
		}

		if toBeDeleted != nil {
			for _, e := range toBeDeleted {
				if err := cs.Events.Delete(data.CalendarID, e.Id).Do(); err != nil {
					log.Error("failed to reset event", zap.Error(err), zap.Any("event", e), zap.Int("year", y), zap.Int("month", m))
				}
			}
			log.Info("deleted stale events", zap.Int("count", len(toBeDeleted)))
		}

		y, m = nextY, nextM
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
