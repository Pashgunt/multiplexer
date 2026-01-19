package config

import (
	"errors"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const (
	KeyTransports = "transports"
	KeyOptions    = "options"
	Filenames     = ".env"
	EmptyEnvName  = ""
)

type TransportOption = map[string]interface{}

type Environment struct {
}

func (env *Environment) Init() error {
	if err := godotenv.Load(Filenames); err != nil {
		return errors.New("error loading .env file")
	}

	return nil
}

func (env *Environment) Replace(data map[string]interface{}) error {
	options, isset := data[KeyTransports].(TransportOption)[KeyOptions]

	if !isset {
		return errors.New("options not found in data")
	}

	for _, value := range options.(TransportOption) {
		transportOptionValue := value.(TransportOption)

		for envKey, envValue := range transportOptionValue {
			envName := env.extractEnvName(envValue.(string))

			if envName == EmptyEnvName {
				continue
			}

			transportOptionValue[envKey] = env.get(envName)
		}
	}

	return nil
}

func (env *Environment) extractEnvName(envString string) string {
	re := regexp.MustCompile(`%env\(([^)]+)\)%`)

	matches := re.FindStringSubmatch(envString)

	if len(matches) > 1 {
		return matches[1]
	}

	return EmptyEnvName
}

func (env *Environment) get(envName string) string {
	return os.Getenv(envName)
}
