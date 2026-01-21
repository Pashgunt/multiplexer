package logging

import (
	"log/slog"
	"os"
	"transport/pkg/utils/backoff"
)

type KafkaConnectionLogger struct {
	logger *slog.Logger
}

func NewKafkaConnectionLogger(level slog.Level) KafkaConnectionLogger {
	return KafkaConnectionLogger{
		logger: slog.
			New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).
			WithGroup(backoff.GroupNameKafkaConnectionLogger),
	}
}

func (logger *KafkaConnectionLogger) Info(entity KafkaConnectionLogEntity) {
	logger.logger.Info(
		entity.Message,
		"broker", entity.Broker,
	)
}
