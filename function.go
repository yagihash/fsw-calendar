package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"

	"github.com/yagihash/fsw-calendar/calendar"
	"github.com/yagihash/fsw-calendar/config"
	"github.com/yagihash/fsw-calendar/event"
	"github.com/yagihash/fsw-calendar/fetcher"
	"github.com/yagihash/fsw-calendar/logger"
	"github.com/yagihash/fsw-calendar/notify/slack"
)

func Register(ctx context.Context, message *pubsub.Message) (err error) {
	c, err := config.Load()
	if err != nil {
		return err
	}

	log, err := logger.New(c.LogLevel)
	if err != nil {
		return err
	}

	defer func() { _ = log.Sync() }()

	log.Debug("logger is ready")

	notify := slack.New(c.Webhook)
	defer func() {
		if err != nil {
			_ = notify.Warn(context.Background(), fmt.Sprintf("error: %s", err.Error()))
		}
	}()

	var data config.Data
	if err := json.Unmarshal(message.Data, &data); err != nil {
		log.Error("failed to unmarshal message", zap.Error(err))
	}

	log.Info("start processing received data", zap.Any("data", data))

	jst, err := time.LoadLocation(c.Timezone)
	if err != nil {
		log.Error("failed to load timezone", zap.Error(err))
		return err
	}

	y := time.Now().In(jst).Year()
	m := int(time.Now().In(jst).Month())

	f := fetcher.New(c.Hostname, data.Course, data.Class, http.DefaultClient)

	docEvents, err := f.FetchDocEvents(y, m, c.Recurrence)
	if err != nil {
		log.Error("failed to fetch schedule data", zap.Error(err))
	}

	log.Debug("loaded schedules", zap.Any("events", docEvents))

	var fetchedEvents event.Events
	for _, d := range docEvents {
		fetchedEvents = append(fetchedEvents, event.NewFromDocEvent(d))
	}

	cs, err := calendar.New(ctx, data.CalendarID, jst)
	if err != nil {
		log.Error("failed to initialize calendar service", zap.Error(err))
		return err
	}

	existingEvents, err := cs.GetEvents(y, m, c.Recurrence)
	if err != nil {
		log.Error("failed to load existing events from google calendar", zap.Error(err))
		return err
	}

	toBeAdded, toBeDeleted := existingEvents.Diff(fetchedEvents)
	if len(toBeAdded) == 0 && len(toBeDeleted) == 0 {
		log.Debug("no update", zap.Any("existing", existingEvents), zap.Any("fetched", fetchedEvents))
		return nil
	}

	log.Info("need updates", zap.Any("to_be_added", toBeAdded), zap.Any("to_be_deleted", toBeDeleted))

	for _, e := range toBeAdded {
		if e == nil {
			continue
		}

		if err := cs.Insert(e.Event); err != nil {
			log.Error("failed to insert event", zap.Error(err), zap.Any("event", e))
		} else {
			log.Debug("added new event", zap.Any("event", e))
		}
	}

	for _, e := range toBeDeleted {
		if err := cs.Delete(e.Id); err != nil {
			log.Error("failed to delete event", zap.Any("event", e))
		} else {
			log.Debug("deleted stale event", zap.Any("event", e))
		}
	}

	_ = notify.Info(ctx, "updated calendar")

	return nil
}
