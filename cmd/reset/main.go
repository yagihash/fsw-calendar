package main

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/api/calendar/v3"

	"github.com/yagihash/fsw-calendar/config"
	"github.com/yagihash/fsw-calendar/event"
	"github.com/yagihash/fsw-calendar/logger"
)

const (
	ExitOK = iota
	ExitError
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log, err := logger.New()
	if err != nil {
		return ExitError
	}

	defer log.Sync()

	log.Info("logger is ready")

	c, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
		return ExitError
	}

	jst, err := time.LoadLocation(c.Timezone)
	if err != nil {
		log.Fatal("failed to load timezone", zap.Error(err))
		return ExitError
	}

	y := time.Now().In(jst).Year()
	m := int(time.Now().In(jst).Month())

	for i := 0; i < c.Recurrence; i++ {
		ctx := context.Background()
		cs, err := calendar.NewService(ctx)
		if err != nil {
			log.Fatal("failed to access google calendar API", zap.Error(err))
			return ExitError
		}

		nextY, nextM := NextMonth(y, m)

		events, err := cs.Events.List(c.CalendarID).ShowDeleted(false).SingleEvents(true).
			TimeMin(time.Date(y, time.Month(m), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).
			TimeMax(time.Date(nextY, time.Month(nextM), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).Do()
		if err != nil {
			log.Error("err", zap.Error(err))
			return ExitError
		}

		existingEvents := event.Events(events.Items)

		for _, e := range existingEvents {
			if err := cs.Events.Delete(c.CalendarID, e.Id).Do(); err != nil {
				log.Error("failed to reset event", zap.Error(err), zap.Any("event", e), zap.Int("year", y), zap.Int("month", m))
			}
		}

		y, m = nextY, nextM
	}

	return ExitOK
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
