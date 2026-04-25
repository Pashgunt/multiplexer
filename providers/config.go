package providers

import (
	"transport/internal/infrastructure/config"
	appconfig "transport/internal/infrastructure/config/app"
)

func Environment() config.EnvironmentInterface {
	env := config.NewEnvironment()

	if err := env.Init(); err != nil {
		panic(err)
	}

	return env
}

func Config(env config.EnvironmentInterface) appconfig.Config {
	basePath := "configs"
	cfg, err := appconfig.Load(basePath, env)

	if err != nil {
		panic(err)
	}

	return *cfg
}
