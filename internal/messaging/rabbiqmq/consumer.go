package rabbiqmq

import (
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
}

func (c *Consumer) consume(config Config) (*amqp091.Channel, error) {
	connection, err := amqp091.Dial(config.url)

	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()

	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		config.name,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return channel, err
}
