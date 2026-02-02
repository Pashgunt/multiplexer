package appconfig

import (
	"transport/internal/application/observability/logging"
	"transport/internal/infrastructure/config"
	"transport/internal/infrastructure/config/types"
)

type Config struct {
	Logger      logging.AdapterInterface
	Environment config.EnvironmentInterface
	Config      types.Config
}
