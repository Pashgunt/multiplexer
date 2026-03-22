package types

type KafkaOption struct {
	Brokers string `yaml:"broker.list"`
	GroupID string `yaml:"group.id"`
}
