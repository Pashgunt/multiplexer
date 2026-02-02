package kafka

import (
	"context"
	"transport/internal/application/observability/logging"

	"github.com/segmentio/kafka-go"
)

type ConnectionInterface interface {
	ConnectionGetterInterface
	ConnectionSetterInterface
	Close() error
}

type ConnectionGetterInterface interface {
	Consumer() ConsumerInterface
	Config() Config
}

type ConnectionSetterInterface interface {
	SetConsumer(consumer ConsumerInterface)
}

type Connection struct {
	uuid       string
	connection *kafka.Conn
	config     Config
	consumer   ConsumerInterface
	logger     logging.LoggerInterface
}

func (connection *Connection) SetConsumer(consumer ConsumerInterface) {
	connection.consumer = consumer
}

func (connection *Connection) Config() Config {
	return connection.config
}

func (connection *Connection) Consumer() ConsumerInterface {
	return connection.consumer
}

func (connection *Connection) Close() error {
	return connection.connection.Close()
}

func NewConnection(ctx context.Context, config Config, logger logging.LoggerInterface) (ConnectionInterface, error) {
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
