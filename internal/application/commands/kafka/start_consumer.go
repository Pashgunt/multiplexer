package kafkacommand

import (
	"transport/internal/messaging/kafka"
	"transport/pkg/logging"
)

func StartConsumers(config kafka.Config, logger logging.LoggerInterface) kafka.ConsumerInterface {
	return kafka.NewConsumer(config, logger)
}
