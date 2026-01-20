package kafka

import (
	"transport/internal/infrastructure/config/types"

	"github.com/segmentio/kafka-go"
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

func (adapter *Adapter) Consume() []*kafka.Reader {

	var readers []*kafka.Reader

	if len(adapter.configs) == 0 {
		return readers
	}

	for _, config := range adapter.configs {
		readers = append(readers, adapter.consumer.consume(config))
	}

	return readers
}
