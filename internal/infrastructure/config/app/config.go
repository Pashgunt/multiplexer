package appconfig

import (
	"transport/internal/infrastructure/config"
	"transport/internal/infrastructure/config/types"
	"transport/internal/infrastructure/db"
	"transport/pkg/logging"
)

type Config struct {
	Logger      logging.AdapterInterface
	Environment config.EnvironmentInterface
	Config      types.Config
	PgSQL       db.IDB
}
