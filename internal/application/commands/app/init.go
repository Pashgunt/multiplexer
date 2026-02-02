package appcommand

import (
	"transport/internal/application/observability/logging"
	"transport/internal/infrastructure/config"
	appconfig "transport/internal/infrastructure/config/app"
	"transport/internal/infrastructure/config/types"
	"transport/pkg/utils/backoff"
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
	kernel.config.Logger = kernel.initLogger()
	kernel.config.Config = kernel.initTransportConfig()

	return kernel
}

func (kernel *Kernel) initEnvironment() config.EnvironmentInterface {
	env, err := config.NewEnvironment()

	if err != nil {
		panic(err)
	}

	return env
}

func (kernel *Kernel) initTransportConfig() types.Config {
	cfg, err := config.NewLoader(config.NewValidator(), kernel.config.Environment).Load(backoff.ConfigPath)

	if err != nil {
		panic(err)
	}

	return *cfg
}

func (kernel *Kernel) initLogger() logging.AdapterInterface {
	logger := logging.NewAdapter()
	logger.Init([]backoff.LoggerType{backoff.KafkaLogger, backoff.AppLogger})

	return logger
}
