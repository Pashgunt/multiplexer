package kafka

import (
	"strings"
	"transport/internal/infrastructure/config/types"
)

const (
	sepBrokers = ","
)

func convert(config types.Config) []Config {
	if len(config.Topics) == 0 {
		return []Config{}
	}

	kafkaConfig := make(map[string]Config)

	for _, topic := range config.Topics {
		if topic.Options == nil || topic.Options.Kafka == nil {
			continue
		}

		brokers := strings.Split(topic.Options.Kafka.Brokers, sepBrokers)

		if len(brokers) == 0 {
			continue
		}

		for _, broker := range brokers {
			key := broker + topic.Options.Kafka.GroupId
			kafkaBrokerConfig, isset := kafkaConfig[key]

			if !isset {
				kafkaConfig[key] = Config{
					Broker:  broker,
					GroupID: topic.Options.Kafka.GroupId,
					Topics:  topic.ConsumerTopics,
				}

				continue
			}

			kafkaBrokerConfig.Topics = append(kafkaBrokerConfig.Topics, topic.ConsumerTopics...)
		}
	}

	if len(kafkaConfig) == 0 {
		return []Config{}
	}

	kafkaConfigList := make([]Config, 0, len(kafkaConfig))

	for _, kafkaConfigItem := range kafkaConfig {
		kafkaConfigList = append(kafkaConfigList, kafkaConfigItem)
	}

	return kafkaConfigList
}
