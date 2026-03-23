package db

import "transport/internal/infrastructure/config"

const (
	databaseSourceName = "PG_DATABASE_URL"
)

type Params struct {
	DatabaseSourceName string
}

func NewParams(env config.EnvironmentInterface) Params {
	return Params{DatabaseSourceName: env.Get(databaseSourceName)}
}
