package logging

import (
	"transport/pkg/utils/backoff"
)

type LoggerFactory interface {
	CreateLogger(loggerType backoff.LoggerType, level backoff.LoggerLevel) LoggerInterface
}

type DefaultLoggerFactory struct{}

func NewDefaultLoggerFactory() *DefaultLoggerFactory {
	return &DefaultLoggerFactory{}
}

func (factory *DefaultLoggerFactory) CreateLogger(
	loggerType backoff.LoggerType,
	level backoff.LoggerLevel,
) LoggerInterface {
	switch loggerType {
	case backoff.KafkaLogger:
		return NewKafkaConnectionLogger(level.GetSlogLevel())
	case backoff.AppLogger:
		return NewAppLogger(level.GetSlogLevel())
	case backoff.APILogger:
		return NewAPILogger(level.GetSlogLevel())
	default:
		return nil
	}
}
