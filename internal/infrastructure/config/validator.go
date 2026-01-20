package config

import (
	"errors"
	"os"
)

type Validator struct {
}

func (validator *Validator) validateFileExists(configPath string) error {
	info, err := os.Stat(configPath)

	if err != nil {
		return err
	}

	if info.IsDir() {
		return errors.New("config path is a directory, not a file")
	}

	if info.Size() == 0 {
		return errors.New("config file is empty")
	}

	return nil
}
