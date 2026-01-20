package config

import (
	"os"
	"transport/internal/infrastructure/config/types"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	validator   *Validator
	environment *Environment
}

func NewLoader(validator *Validator, environment *Environment) *Loader {
	return &Loader{validator: validator, environment: environment}
}

func (loader *Loader) Load(configPath string) (*types.Config, error) {
	if err := loader.validator.validateFileExists(configPath); err != nil {
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

	if err := loader.environment.replace(dataForReplaceEnvironment); err != nil {
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
