package types

type ConfigTransport struct {
	Options *ConfigTransportOptions `yaml:"options"`
	Topics  map[string]*ConfigTopic `yaml:"topics"`
}
