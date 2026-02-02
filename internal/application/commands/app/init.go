package appcommand

import (
	"transport/internal/infrastructure/config"
	"transport/internal/infrastructure/config/types"
	"transport/pkg/utils/backoff"
)

func Init() *types.Config {
	env, err := config.NewEnvironment()

	if err != nil {
		panic(err)
	}

	cfg, err := config.NewLoader(config.NewValidator(), env).Load(backoff.ConfigPath)

	if err != nil {
		panic(err)
	}

	return cfg
}
