package config

import (
	"os"
	"regexp"
	"transport/internal/application/observability/logging"

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

func NewEnvironment() *Environment {
	return &Environment{}
}

func (env *Environment) Init() error {
	if err := godotenv.Load(Filenames); err != nil {
		return logging.NewAppError("Error loading .env file.")
	}

	return nil
}

func (env *Environment) Replace(data map[string]interface{}) error {
	topics, isset := data[KeyTopics]

	if !isset {
		return logging.NewAppError("Topics not found in data.")
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
