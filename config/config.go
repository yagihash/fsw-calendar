package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Timezone   string `envconfig:"TIMEZONE" default:"Asia/Tokyo"`
	CalendarID string `envconfig:"CALENDAR_ID" required:"true"`
	Recurrence int    `envconfig:"RECURRENCE" default:"2"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}
