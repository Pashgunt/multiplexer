package config

import (
	"os"
	"transport/internal/infrastructure/config/types"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	Validator   *Validator
	Environment *Environment
}

func (loader *Loader) Load(configPath string) (*types.Config, error) {
	if err := loader.Validator.ValidateFileExists(configPath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	return loader.decodeAndReplaceEnv(data)
}

func (loader *Loader) decodeAndReplaceEnv(data []byte) (*types.Config, error) {
	var dataForReplaceEnvironment map[string]interface{}

	if err := yaml.Unmarshal(data, &dataForReplaceEnvironment); err != nil {
		return nil, err
	}

	if err := loader.Environment.Replace(dataForReplaceEnvironment); err != nil {
		return nil, err
	}

	yamlWithEnvironment, err := yaml.Marshal(dataForReplaceEnvironment)

	if err != nil {
		return nil, err
	}

	var config types.Config

	if err = yaml.Unmarshal(yamlWithEnvironment, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
