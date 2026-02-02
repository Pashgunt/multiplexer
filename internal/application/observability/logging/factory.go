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
		return NewKafkaConnectionLogger(slog.LevelDebug)
	default:
		return nil
	}
}
