package kafka

import (
	"context"
	"errors"
	"strings"
	"transport/internal/application/observability/logging"
	kafkaconnection "transport/internal/domain/connection"

	"github.com/segmentio/kafka-go"
)

type ConsumerInterface interface {
	Fetch() (kafka.Message, error)
	Commit(messages []kafka.Message, consumerEntity kafkaconnection.Consumer) error
}

type Consumer struct {
	reader  *kafka.Reader
	isReady bool
	logger  logging.LoggerInterface
}

func NewConsumer(config Config, logger logging.LoggerInterface) ConsumerInterface {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{config.Broker},
			GroupTopics: config.Topics,
			GroupID:     config.GroupID,
		}),
		isReady: true,
		logger:  logger,
	}
}

func (consumer *Consumer) Fetch() (kafka.Message, error) {
	brokers := strings.Join(consumer.reader.Config().Brokers, ",")

	consumer.logger.Info(logging.NewKafkaConnectionLogEntity("Waiting a message", brokers))

	if !consumer.isReady {
		err := logging.NewKafkaConsumerNotReady(brokers, "Consumer is not ready, cant start read message")
		consumer.logger.Error(err)

		return kafka.Message{}, err
	}

	message, err := consumer.reader.ReadMessage(context.Background())

	if err != nil {
		consumer.logger.Info(err)

		return message, err
	}

	consumer.logger.Info(logging.NewKafkaConnectionLogEntity("Success get message "+string(message.Value), brokers))

	return message, nil
}

func (consumer *Consumer) Commit(messages []kafka.Message, consumerEntity kafkaconnection.Consumer) error {
	ctxTimeoutCommit, cancel := context.WithTimeout(context.Background(), consumerEntity.Timeout())
	defer cancel()

	if err := consumer.doCommit(messages, ctxTimeoutCommit, consumerEntity.RetryCountCommit()); err != nil {
		consumer.logger.Error(err)

		return err
	}

	return nil
}

func (consumer *Consumer) doCommit(messages []kafka.Message, ctx context.Context, retryCount int) error {
	err := consumer.reader.CommitMessages(ctx, messages...)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}

		if retryCount > 0 {
			return consumer.doCommit(messages, ctx, retryCount-1)
		}

		return logging.NewKafkaCommitError(messages[0].Topic, messages[0].Partition, messages[0].Offset, err.Error())
	}

	return nil
}
