package config

import (
	"errors"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const (
	KeyTopics    = "topics"
	KeyOptions   = "options"
	Filenames    = ".env"
	EmptyEnvName = ""
)

type TransportOption = map[string]interface{}

type EnvironmentInterface interface {
	Init() error
	Get(envName string) string
	Replace(data map[string]interface{}) error
}

type Environment struct {
}

func NewEnvironment() (*Environment, error) {
	env := &Environment{}
	err := env.Init()

	if err != nil {
		return nil, err
	}

	return env, nil
}

func (env *Environment) Init() error {
	if err := godotenv.Load(Filenames); err != nil {
		return errors.New("error loading .env file")
	}

	return nil
}

func (env *Environment) Replace(data map[string]interface{}) error {
	topics, isset := data[KeyTopics]

	if !isset {
		return errors.New("topics not found in data")
	}

	for _, value := range topics.(TransportOption) {
		transportOptionValue := value.(TransportOption)

		for _, optionsParam := range transportOptionValue[KeyOptions].(TransportOption) {
			optionParamValue := optionsParam.(TransportOption)

			for envKey, envValue := range optionParamValue {
				envName := env.extractEnvName(envValue.(string))

				if envName == EmptyEnvName {
					continue
				}

				optionParamValue[envKey] = env.Get(envName)
			}
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

func (env *Environment) Get(envName string) string {
	return os.Getenv(envName)
}
