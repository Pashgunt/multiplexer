package logging

import (
	"log/slog"
	"os"
	"transport/pkg/utils/backoff"
)

type LoggerInterface interface {
	Info(entity interface{})
	Warning(entity interface{})
	Error(entity interface{})
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

func (logger *KafkaConnectionLogger) Info(entity interface{}) {
	kafkaConnectionLoggerEntity, _ := entity.(KafkaConnectionLogEntity)

	logger.logger.Info(
		kafkaConnectionLoggerEntity.Message,
		"broker", kafkaConnectionLoggerEntity.Broker,
	)
}

func (logger *KafkaConnectionLogger) Warning(entity interface{}) {
	kafkaConnectionLoggerEntity, _ := entity.(KafkaConnectionLogEntity)

	logger.logger.Warn(
		kafkaConnectionLoggerEntity.Message,
		"broker", kafkaConnectionLoggerEntity.Broker,
	)
}

func (logger *KafkaConnectionLogger) Error(entity interface{}) {
	kafkaConnectionLoggerEntity, _ := entity.(KafkaConnectionLogEntity)

	logger.logger.Error(
		kafkaConnectionLoggerEntity.Message,
		"broker", kafkaConnectionLoggerEntity.Broker,
	)
}
