package kafka

import (
	"context"
	"strings"
	"transport/internal/application/observability/logging"
	kafkaconnection "transport/internal/domain/connection"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	isReady bool
	logger  logging.KafkaConnectionLogger
}

func NewConsumer(config Config, logger logging.KafkaConnectionLogger) *Consumer {
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

func (consumer *Consumer) Fetch() kafka.Message {
	consumer.logger.Info(logging.KafkaConnectionLogEntity{
		Message: "Waiting a message",
		Broker:  strings.Join(consumer.reader.Config().Brokers, ","),
	})

	if !consumer.isReady {
		consumer.logger.Info(logging.KafkaConnectionLogEntity{
			Message: "Cannot start read message",
			Broker:  strings.Join(consumer.reader.Config().Brokers, ","),
		})

		return kafka.Message{}
	}

	message, err := consumer.reader.ReadMessage(context.Background())

	if err != nil {
		consumer.logger.Info(logging.KafkaConnectionLogEntity{
			Message: "Fetch message with error" + err.Error(),
			Broker:  strings.Join(consumer.reader.Config().Brokers, ","),
		})

		return message
	}

	consumer.logger.Info(logging.KafkaConnectionLogEntity{
		Message: "Success get message " + string(message.Value),
		Broker:  strings.Join(consumer.reader.Config().Brokers, ","),
	})

	return message
}

func (consumer *Consumer) Commit(messages []kafka.Message, consumerEntity kafkaconnection.Consumer) {
	ctxTimeoutCommit, cancel := context.WithTimeout(context.Background(), consumerEntity.Timeout())
	defer cancel()

	select {
	case <-ctxTimeoutCommit.Done():
		consumer.logger.Info(logging.KafkaConnectionLogEntity{
			Message: "Time for commit messages has expired",
			Broker:  strings.Join(consumer.reader.Config().Brokers, ","),
		})
		return
	default:
	}

	if err := consumer.doCommit(messages, ctxTimeoutCommit, consumerEntity.RetryCountCommit()); err != nil {
		consumer.logger.Info(logging.KafkaConnectionLogEntity{
			Message: "Cannot commit message",
			Broker:  strings.Join(consumer.reader.Config().Brokers, ","),
		})
	}
}

func (consumer *Consumer) doCommit(messages []kafka.Message, context context.Context, retryCount int) error {
	err := consumer.reader.CommitMessages(context, messages...)

	if err != nil && retryCount > 0 {
		return consumer.doCommit(messages, context, retryCount-1)
	}

	return nil
}
