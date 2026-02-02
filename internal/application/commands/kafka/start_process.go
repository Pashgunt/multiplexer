package kafkacommand

import (
	appconfig "transport/internal/infrastructure/config/app"
	"transport/internal/messaging/kafka"
	"transport/pkg/utils/backoff"
)

func StartProcess(connections []kafka.ConnectionInterface, config appconfig.Config) {
	for _, connection := range connections {
		go doProcessForConsumer(connection, config)
	}
}

func doProcessForConsumer(connection kafka.ConnectionInterface, config appconfig.Config) {
	connection.SetConsumer(StartConsumers(connection.Config(), config.Logger.GetLogger(backoff.KafkaLogger)))
	ConsumeMessage(connection.Consumer())
}
