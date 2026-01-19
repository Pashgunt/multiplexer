package types

type RabbitMQOption struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Vhost string `yaml:"vhost"`
}
