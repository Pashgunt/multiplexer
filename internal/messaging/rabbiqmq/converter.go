package rabbiqmq

import (
	"strings"
	"transport/internal/infrastructure/config/types"
)

const (
	amqpPrefix          = "amqp://"
	amqpPortDelimiter   = ':'
	amqpPassDelimiter   = '@'
	countDelimiterChars = 3
)

func convert(config types.Config) []Config {
	var amqpConfig []Config

	if len(config.Topics) == 0 {
		return amqpConfig
	}

	for _, topic := range config.Topics {
		if topic.Options == nil || topic.Options.RabbitMQ == nil || len(topic.ConsumerTopics) == 0 {
			continue
		}

		for _, topicName := range topic.ConsumerTopics {
			amqpConfig = append(amqpConfig, Config{
				url:  formatUrl(topic.Options.RabbitMQ),
				name: topicName,
			})
		}
	}

	return amqpConfig
}

func formatUrl(option *types.RabbitMQOption) string {
	if option == nil {
		return ""
	}

	var builder strings.Builder

	builder.Grow(len(amqpPrefix) + len(option.Username) + len(option.Password) +
		len(option.Host) + len(option.Port) + len(option.Vhost) + countDelimiterChars)

	builder.WriteString(amqpPrefix)
	builder.WriteString(option.Username)
	builder.WriteByte(amqpPortDelimiter)
	builder.WriteString(option.Password)
	builder.WriteByte(amqpPassDelimiter)
	builder.WriteString(option.Host)
	builder.WriteByte(amqpPortDelimiter)
	builder.WriteString(option.Port)
	builder.WriteString(option.Vhost)

	return builder.String()
}
