package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/yagihash/fsw-calendar/fetcher/class"
	"github.com/yagihash/fsw-calendar/fetcher/course"
)

func TestData_UnmarshalJSON(t *testing.T) {
	opt := cmp.Transformer("", func(src fmt.Stringer) string {
		return src.String()
	})

	input := `{"calendar_id":"https://example.com/test","course":"rc","class":"t-4"}`

	var got Data
	if err := json.Unmarshal([]byte(input), &got); err != nil {
		t.Errorf("got an unexpected error: %v", err)
	}

	want := Data{
		CalendarID: "https://example.com/test",
		Course:     course.RC,
		Class:      class.T4,
	}

	if diff := cmp.Diff(got, want, opt); diff != "" {
		t.Errorf("got unexpected diff:\n%s", diff)
	}
}

func TestData_MarshalLogObject(t *testing.T) {
	var buf bytes.Buffer

	syncer := zapcore.AddSync(&buf)
	lock := zapcore.Lock(syncer)
	encoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(encoder, lock, zapcore.InfoLevel)
	logger := zap.New(core)

	data := &Data{
		CalendarID: "test@example.com",
		Course:     course.SS,
		Class:      class.SS4,
	}

	logger.Info("test log", zap.Any("data", data))

	var got map[string]interface{}
	if err := json.Unmarshal([]byte(buf.String()), &got); err != nil {
		t.Errorf("failed to unmarshal json log string: %v", err)
	}

	wantData := map[string]interface{}{
		"calendar_id": "test@example.com",
		"course":      "ss",
		"class":       "ss-4",
	}

	if diff := cmp.Diff(got["data"], wantData); diff != "" {
		t.Errorf("got unexpected diff:\n%s", diff)
	}
}
