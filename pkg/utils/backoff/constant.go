package backoff

const (
	Test    = "test"
	Dev     = "dev"
	PreProd = "pre-prod"
	Prod    = "prod"
)

const (
	GroupNameKafkaConnectionLogger = "kafka.connection"
)

const (
	ConfigPath = "./configs/transport.yaml"
)

type LoggerType string

const (
	KafkaLogger LoggerType = "kafka"
	AppLogger   LoggerType = "app"
)
