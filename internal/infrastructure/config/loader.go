package config

import (
	"os"
	"transport/internal/infrastructure/config/types"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	Validator *Validator
}

func (loader *Loader) Load(configPath string) (*types.Config, error) {
	if err := loader.Validator.ValidateFileExists(configPath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	config := types.Config{}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
