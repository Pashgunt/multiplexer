package logging

import (
	"log/slog"
	"os"
	"transport/pkg/utils/backoff"
)

type LoggerInterface interface {
	Info(interface{})
	Warning(interface{})
	Error(error)
}

type KafkaConnectionLogger struct {
	logger *slog.Logger
}

func NewKafkaConnectionLogger(level slog.Level) LoggerInterface {
	return &KafkaConnectionLogger{
		logger: slog.
			New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).
			WithGroup(backoff.GroupNameKafkaConnectionLogger),
	}
}

func (logger *KafkaConnectionLogger) Warning(object interface{}) {
	kafkaConnectionLogEntity := object.(KafkaConnectionLogEntity)

	logger.logger.Warn(kafkaConnectionLogEntity.Message, "broker", kafkaConnectionLogEntity.Broker)
}

func (logger *KafkaConnectionLogger) Info(object interface{}) {
	kafkaConnectionLogEntity := object.(KafkaConnectionLogEntity)

	logger.logger.Info(kafkaConnectionLogEntity.Message, "broker", kafkaConnectionLogEntity.Broker)
}

func (logger *KafkaConnectionLogger) Error(error error) {
	logger.logger.Error(error.Error())
}

type AppLogger struct {
	logger *slog.Logger
}

func NewAppLogger(level slog.Level) LoggerInterface {
	return &AppLogger{
		logger: slog.
			New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).
			WithGroup(backoff.GroupNameAppLogger),
	}
}

func (logger *AppLogger) Warning(object interface{}) {
	appLogEntity := object.(AppLogEntity)

	logger.logger.Warn(appLogEntity.Message)
}

func (logger *AppLogger) Info(object interface{}) {
	appLogEntity := object.(AppLogEntity)

	logger.logger.Info(appLogEntity.Message)
}

func (logger *AppLogger) Error(error error) {
	logger.logger.Error(error.Error())
}
