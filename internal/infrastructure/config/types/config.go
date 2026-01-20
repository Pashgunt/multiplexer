package types

type Config struct {
	Topics map[string]*ConfigTopic `yaml:"topics"`
}
