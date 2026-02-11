package logging

import "fmt"

type KafkaConnectionError struct {
	Broker  string
	Message string
}

func NewKafkaConnectionError(broker, message string) *KafkaConnectionError {
	return &KafkaConnectionError{Broker: broker, Message: message}
}

func (k KafkaConnectionError) Error() string {
	return fmt.Sprintf("%s: Broker: %s", k.Message, k.Broker)
}

type KafkaCommitError struct {
	Topic     string
	Partition int
	Offset    int64
	Message   string
}

func NewKafkaCommitError(topic string, partition int, offset int64, message string) *KafkaCommitError {
	return &KafkaCommitError{Topic: topic, Partition: partition, Offset: offset, Message: message}
}

func (k KafkaCommitError) Error() string {
	return fmt.Sprintf("%s. Topic: %s. Partition: %d. Offset: %d.", k.Message, k.Topic, k.Partition, k.Offset)
}

type KafkaConsumerNotReady struct {
	Broker  string
	Message string
}

func NewKafkaConsumerNotReady(broker string, message string) *KafkaConsumerNotReady {
	return &KafkaConsumerNotReady{Broker: broker, Message: message}
}

func (k KafkaConsumerNotReady) Error() string {
	return fmt.Sprintf("%s. Broker: %s.", k.Message, k.Broker)
}

type AppError struct {
	Message string
}

func NewAppError(message string) *AppError {
	return &AppError{Message: message}
}

func (a AppError) Error() string {
	return a.Message
}
