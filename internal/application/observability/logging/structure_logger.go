package logging

type KafkaConnectionLogEntity struct {
	Message string
	Broker  string
}

func NewKafkaConnectionLogEntity(message string, broker string) KafkaConnectionLogEntity {
	return KafkaConnectionLogEntity{Message: message, Broker: broker}
}

type AppLogEntity struct {
	Message string
}

func NewAppLogEntity(message string) AppLogEntity {
	return AppLogEntity{Message: message}
}
