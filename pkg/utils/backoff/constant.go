package backoff

const (
	GroupNameKafkaConnectionLogger = "kafka.connection"
	GroupNameAppLogger             = "app"
	GroupNameApiLogger             = "api"
)

const (
	ConfigPath = "./configs/transport.yaml"
)

type LoggerType string

const (
	KafkaLogger LoggerType = "kafka"
	AppLogger   LoggerType = "app"
	ApiLogger   LoggerType = "api"
)
