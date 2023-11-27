package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Timezone   string        `envconfig:"TIMEZONE" default:"Asia/Tokyo"`
	Recurrence int           `envconfig:"RECURRENCE" default:"2"`
	LogLevel   zapcore.Level `envconfig:"LOG_LEVEL" default:"INFO"`
	Hostname   string        `envconfig:"HOSTNAME" default:"www.fsw.tv"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}
