package appconfig

import (
	"fmt"
	"transport/internal/infrastructure/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP  HTTPConfig  `yaml:"HTTP"`
	DB    DBConfig    `yaml:"DB"`
	Redis RedisConfig `yaml:"Redis"`
}

type HTTPConfig struct {
	Port              string `yaml:"Port"`
	Host              string `yaml:"Host"`
	ReadTimeout       int    `yaml:"ReadTimeout"`
	WriteTimeout      int    `yaml:"WriteTimeout"`
	ReadHeaderTimeout int    `yaml:"ReadHeaderTimeout"`
}

type DBConfig struct {
	DatabaseSourceName string `yaml:"DatabaseSourceName"`
}

type RedisConfig struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

func Load(baseConfigPath string, env config.EnvironmentInterface) (*Config, error) {
	appEnv := env.Get("APP_ENV")

	if appEnv == "" {
		appEnv = "local"
	}

	path := baseConfigPath + "/" + appEnv + ".config.yaml"

	fmt.Println(path)
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
