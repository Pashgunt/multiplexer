package kafka

import (
	"strings"
	"transport/internal/infrastructure/config/types"
)

const (
	SepBrokers = ","
)

func convert(config types.Config) []Config {
	var kafkaConfig []Config

	if len(config.Topics) == 0 {
		return kafkaConfig
	}

	for _, topic := range config.Topics {
		if topic.Options == nil || topic.Options.Kafka == nil {
			continue
		}

		kafkaConfig = append(kafkaConfig, Config{
			Brokers: strings.Split(topic.Options.Kafka.Brokers, SepBrokers),
			GroupID: topic.Options.Kafka.GroupId,
			Topics:  topic.ConsumerTopics,
		})
	}

	return kafkaConfig
}
