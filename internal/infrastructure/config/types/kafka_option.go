package types

type KafkaOption struct {
	Dsn     string `yaml:"dsn"`
	Brokers string `yaml:"broker.list"`
	GroupId string `yaml:"group.id"`
}
