package kafka

import (
	"context"
	"transport/internal/application/observability/logging"

	"github.com/segmentio/kafka-go"
)

type Connection struct {
	uuid       string
	connection *kafka.Conn
	config     Config
	consumer   *Consumer
	logger     logging.KafkaConnectionLogger
}

func (connection *Connection) SetConsumer(consumer *Consumer) {
	connection.consumer = consumer
}

func (connection *Connection) Config() Config {
	return connection.config
}

func (connection *Connection) Consumer() *Consumer {
	return connection.consumer
}

func (connection *Connection) Close() error {
	return connection.connection.Close()
}

func NewConnection(ctx context.Context, config Config, logger logging.KafkaConnectionLogger) (*Connection, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	connection, err := kafka.Dial("tcp", config.Broker)

	if err != nil {
		return nil, err
	}

	if _, err = connection.Brokers(); err != nil {
		err = connection.Close()

		if err != nil {
			return nil, err
		}

		return nil, err
	}

	return &Connection{
		connection: connection,
		config:     config,
		logger:     logger,
	}, nil
}
