package rabbiqmq

import (
	"transport/internal/infrastructure/config/types"

	"github.com/rabbitmq/amqp091-go"
)

type Adapter struct {
	consumer *Consumer
	configs  []Config
}

func NewAdapter(appConfig types.Config) *Adapter {
	return &Adapter{
		consumer: &Consumer{},
		configs:  convert(appConfig),
	}
}

func (adapter *Adapter) Consume() ([]*amqp091.Channel, error) {

	var readers []*amqp091.Channel

	if len(adapter.configs) == 0 {
		return readers, nil
	}

	for _, config := range adapter.configs {
		channel, err := adapter.consumer.consume(config)

		if err != nil {
			return nil, err
		}

		readers = append(readers, channel)
	}

	return readers, nil
}
