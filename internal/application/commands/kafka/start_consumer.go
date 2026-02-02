package kafkacommand

import (
	"transport/internal/application/observability/logging"
	"transport/internal/messaging/kafka"
)

func StartConsumers(config kafka.Config, logger logging.LoggerInterface) kafka.ConsumerInterface {
	return kafka.NewConsumer(config, logger)
}
