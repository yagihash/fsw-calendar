package config

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

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
		fmt.Errorf("got unexpected diff:\n%s", diff)
	}
}
