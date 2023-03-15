package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/api/calendar/v3"

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

	url, calendarID, year, month, err := LoadEnv()
	if err != nil {
		log.Fatal("lack of required env var", zap.Error(err))
		return ExitError
	}

	y, _ := strconv.Atoi(year)
	m, _ := strconv.Atoi(month)

	nextY, nextM := NextMonth(y, m)

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal("failed to load timezone", zap.Error(err))
		return ExitError
	}

	fetchedEvents, err := event.Fetch(url)
	if err != nil {
		log.Fatal("failed to fetch schedule data", zap.Error(err), zap.String("url", url))
		return ExitError
	}

	log.Info("loaded schedules", zap.String("url", url))

	cs, err := calendar.NewService(context.TODO())
	if err != nil {
		log.Fatal("failed to access google calendar API", zap.Error(err))
		return ExitError
	}

	events, err := cs.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).
		TimeMin(time.Date(y, time.Month(m), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).
		TimeMax(time.Date(nextY, time.Month(nextM), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).Do()
	if err != nil {
		log.Error("err", zap.Error(err))
		return ExitError
	}

	existingEvents := event.NewEvents(events.Items)

	log.Info("events length", zap.Int("length", len(fetchedEvents)))
	fetchedEvents = fetchedEvents.Unique()
	log.Info("events length", zap.Int("length", len(fetchedEvents)))

	toBeAdded, toBeDeleted := existingEvents.Diff(fetchedEvents)
	if toBeAdded == nil && toBeDeleted == nil {
		log.Info("no update", zap.Int("year", y), zap.Int("month", m))
	}

	if toBeAdded != nil {
		for _, e := range toBeAdded {
			if e == nil {
				continue
			}

			_, err := cs.Events.Insert(calendarID, e).Do()
			if err != nil {
				log.Error("failed to insert event", zap.Error(err), zap.Any("event", e), zap.Int("year", y), zap.Int("month", m))
			}
		}
		log.Info("added new events", zap.Int("count", len(toBeAdded)))
	}

	if toBeDeleted != nil {
		for _, e := range toBeDeleted {
			if err := cs.Events.Delete(calendarID, e.Id).Do(); err != nil {
				log.Error("failed to reset event", zap.Error(err), zap.Any("event", e), zap.Int("year", y), zap.Int("month", m))
			}
		}
		log.Info("deleted stale events", zap.Int("count", len(toBeDeleted)))
	}

	return ExitOK
}

func LoadEnv() (url, calendarID, year, month string, err error) {
	url, ok := os.LookupEnv("URL")
	if !ok {
		return "", "", "", "", fmt.Errorf("no url is specified")
	}

	calendarID, ok = os.LookupEnv("CALENDAR_ID")
	if !ok {
		return "", "", "", "", fmt.Errorf("no calendar_id is specified")
	}

	year, ok = os.LookupEnv("YEAR")
	if !ok {
		return "", "", "", "", fmt.Errorf("no year is specified")
	}

	month, ok = os.LookupEnv("MONTH")
	if !ok {
		return "", "", "", "", fmt.Errorf("no month is specified")
	}

	return
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
