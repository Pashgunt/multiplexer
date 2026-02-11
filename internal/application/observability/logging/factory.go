package logging

import (
	"log/slog"
	"transport/pkg/utils/backoff"
)

type LoggerFactory interface {
	CreateLogger(loggerType backoff.LoggerType) LoggerInterface
}

type defaultLoggerFactory struct{}

func (factory *defaultLoggerFactory) CreateLogger(loggerType backoff.LoggerType) LoggerInterface {
	switch loggerType {
	case backoff.KafkaLogger:
		return NewKafkaConnectionLogger(slog.LevelDebug) //todo set debug level with config file
	case backoff.AppLogger:
		return NewAppLogger(slog.LevelDebug) //todo set debug level with config file
	default:
		return nil
	}
}
