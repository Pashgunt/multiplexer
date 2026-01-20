package kafka

type Config struct {
	Brokers []string
	Topics  []string
	GroupID string
}
