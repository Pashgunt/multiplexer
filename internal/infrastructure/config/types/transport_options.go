package types

type ConfigTransportOptions struct {
	Kafka       *KafkaOption       `yaml:"kafka,omitempty"`
	RabbitMQ    *RabbitMQOption    `yaml:"rabbitmq,omitempty"`
	RedisStream *RedisStreamOption `yaml:"redis_stream,omitempty"`
}
