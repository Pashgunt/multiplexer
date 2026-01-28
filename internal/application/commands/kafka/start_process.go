package kafkacommand

import (
	"transport/internal/application/observability/logging"
	"transport/internal/messaging/kafka"
)

func StartProcess(connections []*kafka.Connection, logger logging.KafkaConnectionLogger) {
	for _, connection := range connections {
		go doProcessForConsumer(connection, logger)
	}
}

func doProcessForConsumer(connection *kafka.Connection, logger logging.KafkaConnectionLogger) {
	connection.SetConsumer(StartConsumers(connection.Config(), logger))
	ConsumeMessage(connection.Consumer())
}
