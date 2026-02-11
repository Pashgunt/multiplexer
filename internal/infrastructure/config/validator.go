package config

import (
	"os"
	"transport/internal/application/observability/logging"
)

type ValidatorTransportStructInterface interface {
	ValidateFileExists(configPath string) error
}

type Validator struct {
}

func NewValidator() ValidatorTransportStructInterface {
	return &Validator{}
}

func (validator *Validator) ValidateFileExists(configPath string) error {
	info, err := os.Stat(configPath)

	if err != nil {
		return err
	}

	if info.IsDir() {
		return logging.NewAppError("Config path is a directory, not a file.")
	}

	if info.Size() == 0 {
		return logging.NewAppError("Config file is empty.")
	}

	return nil
}
