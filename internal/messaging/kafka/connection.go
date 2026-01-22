package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Connection struct {
	uuid       string
	connection *kafka.Conn
	config     Config
	consumers  []*Consumer
}

func (connection *Connection) Close() error {
	return connection.connection.Close()
}

func NewConnection(ctx context.Context, config Config) (*Connection, error) {
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
	}, nil
}
