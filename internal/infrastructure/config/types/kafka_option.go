package types

type KafkaOption struct {
	Brokers string `yaml:"broker.list"`
	GroupId string `yaml:"group.id"`
}
