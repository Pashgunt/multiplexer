package backoff

import (
	"log/slog"
)

const (
	GroupNameKafkaConnectionLogger = "kafka.connection"
	GroupNameAppLogger             = "app"
	GroupNameAPILogger             = "api"
)

const (
	ConfigPath = "./configs/transport.yaml"
)

type LoggerType string
type LoggerLevel string

func (l LoggerLevel) GetSlogLevel() slog.Level {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

const (
	Info  LoggerLevel = "INFO"
	Warn  LoggerLevel = "WARN"
	Error LoggerLevel = "ERROR"
	Debug LoggerLevel = "DEBUG"
)

const (
	KafkaLogger LoggerType = "kafka"
	AppLogger   LoggerType = "app"
	APILogger   LoggerType = "api"
)

const (
	EnvKafkaDebugLevelKey = "KAFKA_DEBUG_LEVEL"
	EnvAppDebugLevelKey   = "APP_DEBUG_LEVEL"
	EnvAPIDebugLevelKey   = "API_DEBUG_LEVEL"
)
