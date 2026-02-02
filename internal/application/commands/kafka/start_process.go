package kafkacommand

import (
	"transport/internal/application/observability/logging"
	"transport/internal/messaging/kafka"
)

func StartProcess(connections []kafka.ConnectionInterface, logger logging.LoggerInterface) {
	for _, connection := range connections {
		go doProcessForConsumer(connection, logger)
	}
}

func doProcessForConsumer(connection kafka.ConnectionInterface, logger logging.LoggerInterface) {
	connection.SetConsumer(StartConsumers(connection.Config(), logger))
	ConsumeMessage(connection.Consumer())
}
