package types

type ConfigTopic struct {
	Options        *ConfigTransportOptions `yaml:"options"`
	ConsumerTopics []string                `yaml:"consumer_topics"`
}
