package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	EnvTimezone   = "TIMEZONE"
	EnvCalendarID = "CALENDAR_ID"
	EnvRecurrence = "RECURRENCE"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]string
		want        *Config
		expectError bool
	}{
		{
			name: "SuccessWithSpecifiedValues",
			input: map[string]string{
				EnvTimezone:   "TEST_TIMEZONE",
				EnvCalendarID: "TEST_CALENDAR_ID",
				EnvRecurrence: "10",
			},
			want: &Config{
				Timezone:   "TEST_TIMEZONE",
				CalendarID: "TEST_CALENDAR_ID",
				Recurrence: 10,
			},
			expectError: false,
		},
		{
			name: "SuccessWithDefaultTimezone",
			input: map[string]string{
				EnvCalendarID: "TEST_CALENDAR_ID",
				EnvRecurrence: "10",
			},
			want: &Config{
				Timezone:   "Asia/Tokyo",
				CalendarID: "TEST_CALENDAR_ID",
				Recurrence: 10,
			},
			expectError: false,
		},
		{
			name: "SuccessWithDefaultRecurrence",
			input: map[string]string{
				EnvTimezone:   "TEST_TIMEZONE",
				EnvCalendarID: "TEST_CALENDAR_ID",
			},
			want: &Config{
				Timezone:   "TEST_TIMEZONE",
				CalendarID: "TEST_CALENDAR_ID",
				Recurrence: 2,
			},
			expectError: false,
		},
		{
			name:        "FailWithMissingRequiredValue",
			input:       map[string]string{},
			want:        nil,
			expectError: true,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			if err := SetEnv(c.input); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := UnsetEnv(c.input); err != nil {
					t.Fatal(err)
				}
			}()

			got, err := Load()

			if err == nil && c.expectError {
				t.Errorf("got an unexpected error: %v", err)
				return
			}

			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("got an unexpected diff:\n%s", diff)
			}
		})
	}
}

func SetEnv(env map[string]string) error {
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnsetEnv(env map[string]string) error {
	for k, _ := range env {
		err := os.Unsetenv(k)
		if err != nil {
			return err
		}
	}

	return nil
}
