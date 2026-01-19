package types

type RedisStreamOption struct {
	Address string `yaml:"address"`
	DB      int    `yaml:"db"`
}
