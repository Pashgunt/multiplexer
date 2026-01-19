package types

type RabbitMQOption struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Vhost string `yaml:"vhost"`
}
