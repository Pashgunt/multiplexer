package kafkacommand

import (
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/messaging/kafka"

	kafkago "github.com/segmentio/kafka-go"
)

func ConsumeMessage(consumer kafka.ConsumerInterface) {
	for {
		message, err := consumer.Fetch()

		if err != nil {
			continue
		}

		_ = consumer.Commit([]kafkago.Message{message}, kafkaconnection.DefaultConsumer())
	}
}
