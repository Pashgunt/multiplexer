package backoff

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
