package config

import (
	"os"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/google/go-cmp/cmp"
)

const (
	EnvTimezone   = "TIMEZONE"
	EnvRecurrence = "RECURRENCE"
	EnvLogLevel   = "LOG_LEVEL"
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
				EnvRecurrence: "10",
				EnvLogLevel:   "DEBUG",
			},
			want: &Config{
				Timezone:   "TEST_TIMEZONE",
				Recurrence: 10,
				LogLevel:   zapcore.DebugLevel,
			},
			expectError: false,
		},
		{
			name: "SuccessWithDefaultTimezone",
			input: map[string]string{
				EnvRecurrence: "10",
				EnvLogLevel:   "FATAL",
			},
			want: &Config{
				Timezone:   "Asia/Tokyo",
				Recurrence: 10,
				LogLevel:   zapcore.FatalLevel,
			},
			expectError: false,
		},
		{
			name: "SuccessWithDefaultRecurrence",
			input: map[string]string{
				EnvTimezone: "TEST_TIMEZONE",
				EnvLogLevel: "FATAL",
			},
			want: &Config{
				Timezone:   "TEST_TIMEZONE",
				Recurrence: 2,
				LogLevel:   zapcore.FatalLevel,
			},
			expectError: false,
		},
		{
			name: "SuccessWithDefaultLogLevel",
			input: map[string]string{
				EnvTimezone:   "TEST_TIMEZONE",
				EnvRecurrence: "10",
			},
			want: &Config{
				Timezone:   "TEST_TIMEZONE",
				Recurrence: 10,
			},
			expectError: false,
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

			if c.expectError && err == nil {
				t.Error("expected error but got nil")
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
