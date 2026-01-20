package types

type RabbitMQOption struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Vhost    string `yaml:"vhost"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
