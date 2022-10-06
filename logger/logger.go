package logger

import (
	"cloud.google.com/go/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevel = map[zapcore.Level]string{
	zapcore.DebugLevel:  logging.Debug.String(),
	zapcore.InfoLevel:   logging.Info.String(),
	zapcore.WarnLevel:   logging.Warning.String(),
	zapcore.ErrorLevel:  logging.Error.String(),
	zapcore.DPanicLevel: logging.Critical.String(),
	zapcore.PanicLevel:  logging.Critical.String(),
	zapcore.FatalLevel:  logging.Critical.String(),
}

func EncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString(logging.Debug.String())
	case zapcore.InfoLevel:
		enc.AppendString(logging.Info.String())
	case zapcore.WarnLevel:
		enc.AppendString(logging.Warning.String())
	case zapcore.ErrorLevel:
		enc.AppendString(logging.Error.String())
	default:
		enc.AppendString(logging.Critical.String())
	}
}

func New() (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "severity"
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeLevel = EncodeLevel
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.ErrorLevel)
	cfg.EncoderConfig = encoderConfig

	return cfg.Build()
}
