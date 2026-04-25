package providers

import (
	"database/sql"
	appconfig "transport/internal/infrastructure/config/app"

	"github.com/pressly/goose/v3"
)

func DB(cfg appconfig.Config) *sql.DB {
	db, err := sql.Open(string(goose.DialectPostgres), cfg.DB.DatabaseSourceName)

	if err != nil {
		panic(err)
	}

	return db
}
