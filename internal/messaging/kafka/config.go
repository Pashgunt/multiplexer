package kafka

type Config struct {
	Broker  string
	Topics  []string
	GroupID string
}
