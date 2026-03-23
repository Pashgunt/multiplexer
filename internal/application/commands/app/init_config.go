package appcommand

import (
	"transport/internal/infrastructure/config"
	appconfig "transport/internal/infrastructure/config/app"
)

const (
	envNamePgDatabaseURL = "PG_DATABASE_URL"
)

type KernelInterface interface {
	KernelGetterInterface
	Init() KernelInterface
}

type KernelGetterInterface interface {
	Config() appconfig.Config
}

type Kernel struct {
	config appconfig.Config
}

func NewKernel() KernelInterface {
	return &Kernel{config: appconfig.Config{}}
}

func (kernel *Kernel) Config() appconfig.Config {
	return kernel.config
}

func (kernel *Kernel) Init() KernelInterface {
	kernel.config.Environment = kernel.initEnvironment()

	return kernel
}

func (kernel *Kernel) initEnvironment() config.EnvironmentInterface {
	env := config.NewEnvironment()

	if err := env.Init(); err != nil {
		panic(err)
	}

	return env
}
