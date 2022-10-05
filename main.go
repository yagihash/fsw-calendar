package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/yagihash/fsw-calendar/frame"

	"google.golang.org/api/calendar/v3"

	"go.uber.org/zap"

	"github.com/yagihash/fsw-calendar/config"
)

const (
	ExitOK = iota
	ExitError

	TemplateURL = "https://www.fsw.tv/driving/sports/ss/ss-4/%d/%02d.html"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println(err)
		return ExitError
	}

	defer logger.Sync()

	logger.Info("logger is ready")

	c, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
		return ExitError
	}

	jst, err := time.LoadLocation(c.Timezone)
	if err != nil {
		logger.Fatal("failed to load timezone", zap.Error(err))
		return ExitError
	}

	y := time.Now().In(jst).Year()
	m := int(time.Now().In(jst).Month())

	for i := 0; i < c.Recurrence; i++ {
		url := fmt.Sprintf(TemplateURL, y, m)

		logger.Info("start to load schedules", zap.String("url", url))

		fetchedFrames, err := frame.Fetch(url)
		if err != nil {
			logger.Fatal("failed to fetch schedule data", zap.Error(err))
			return ExitError
		}

		ctx := context.Background()
		cs, err := calendar.NewService(ctx) //, option.WithCredentialsFile("yagihash-892cb09a93a9.json"))
		if err != nil {
			logger.Fatal("failed to access google calendar API", zap.Error(err))
			return ExitError
		}

		nextY, nextM := NextMonth(y, m)

		events, err := cs.Events.List(c.CalendarID).ShowDeleted(false).SingleEvents(true).
			TimeMin(time.Date(y, time.Month(m), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).
			TimeMax(time.Date(nextY, time.Month(nextM), 1, 0, 0, 0, 0, jst).Format(time.RFC3339)).Do()
		if err != nil {
			logger.Error("err", zap.Error(err))
		}

		var existingFrames []*frame.Frame
		for _, e := range events.Items {
			existingFrames = append(existingFrames, frame.NewFromEvent(e))
		}

		if reflect.DeepEqual(existingFrames, fetchedFrames) {
			logger.Info("no update", zap.Int("year", y), zap.Int("month", m))
		} else {
			logger.Info("start replacing or adding events", zap.Int("year", y), zap.Int("month", m))

			for _, e := range events.Items {
				if err := cs.Events.Delete(c.CalendarID, e.Id).Do(); err != nil {
					logger.Error("failed to delete event", zap.Error(err), zap.Any("event", e))
				}
			}

			for _, f := range fetchedFrames {
				if f == nil {
					continue
				}

				_, err := cs.Events.Insert(c.CalendarID, f.Event()).Do()
				if err != nil {
					logger.Error("failed to insert event", zap.Error(err), zap.Any("event", f.Event()))
					return ExitError
				}
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
