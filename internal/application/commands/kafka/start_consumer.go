package kafkacommand

import (
	"transport/internal/application/observability/logging"
	"transport/internal/messaging/kafka"
)

func StartConsumers(config kafka.Config, logger logging.Logger) *kafka.Consumer {
	return kafka.NewConsumer(config, logger)
}
