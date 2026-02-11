package logging

type KafkaConnectionLogEntity struct {
	Message string
	Broker  string
}

func NewKafkaConnectionLogEntity(message string, broker string) KafkaConnectionLogEntity {
	return KafkaConnectionLogEntity{Message: message, Broker: broker}
}
