package types

type ConfigOption struct {
	Dsn     string `yaml:"dsn"`
	Brokers string `yaml:"broker.list"`
	GroupId string `yaml:"group.id"`
}
