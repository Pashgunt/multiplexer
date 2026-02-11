package appconfig

import (
	"transport/internal/application/observability/logging"
	"transport/internal/infrastructure/config"
	"transport/internal/infrastructure/config/types"
	"transport/internal/infrastructure/db"
)

type Config struct {
	Logger      logging.AdapterInterface
	Environment config.EnvironmentInterface
	Config      types.Config
	PgSql       db.DBInterface
}
