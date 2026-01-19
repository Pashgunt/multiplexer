package types

type ConfigTopic struct {
	Type           string   `yaml:"type"`
	ConsumerTopics []string `yaml:"consumer_topics"`
}
