package kafka

import (
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
}

func (consumer *Consumer) consume(config Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Brokers,
		GroupTopics: config.Topics,
		GroupID:     config.GroupID,
	})
}
