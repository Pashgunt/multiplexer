package main

import (
	"flag"
	"fmt"
	"transport/api/database/migrations"
	appcommand "transport/internal/application/commands/app"
	"transport/internal/application/observability/logging"
	"transport/pkg/utils/backoff"

	_ "github.com/lib/pq"
)

const (
	gooseResetOption  = "reset"
	gooseDownOption   = "down"
	gooseStatusOption = "status"
)

func main() {
	reset := flag.Bool(gooseResetOption, false, "Reset all migrations.")
	down := flag.Bool(gooseDownOption, false, "Down last migration.")
	status := flag.Bool(gooseStatusOption, false, "Get status migration.")
	flag.Parse()

	app := appcommand.NewKernel().Init().Config()

	defer app.PgSql.Close()

	migrator := migrations.NewMigrator(app.PgSql.Db())

	if err := migrator.Setup(); err != nil {
		app.
			Logger.
			GetLogger(backoff.AppLogger).
			Error(logging.NewAppError(fmt.Sprintf("Failed to setup migrator: %s", err.Error())))
	}

	version, _ := migrator.Version()
	app.
		Logger.
		GetLogger(backoff.AppLogger).
		Info(logging.NewAppLogEntity(fmt.Sprintf("Current migration version: %d", version)))

	switch {
	case *reset:
		if err := migrator.Reset(); err != nil {
			app.
				Logger.
				GetLogger(backoff.AppLogger).
				Error(logging.NewAppError(fmt.Sprintf("Failed to reset:: %s", err.Error())))
		}
	case *down:
		if err := migrator.Down(); err != nil {
			app.
				Logger.
				GetLogger(backoff.AppLogger).
				Error(logging.NewAppError(fmt.Sprintf("Failed to rollback:: %s", err.Error())))
		}
	case *status:
		if err := migrator.Status(); err != nil {
			app.
				Logger.
				GetLogger(backoff.AppLogger).
				Error(logging.NewAppError(fmt.Sprintf("Failed to get status:: %s", err.Error())))
		}
	default:
		if err := migrator.Up(); err != nil {
			app.
				Logger.
				GetLogger(backoff.AppLogger).
				Error(logging.NewAppError(fmt.Sprintf("Failed to apply migrations:: %s", err.Error())))
		}
	}

	newVersion, _ := migrator.Version()
	app.
		Logger.
		GetLogger(backoff.AppLogger).
		Info(logging.NewAppLogEntity(fmt.Sprintf("New migration version: %d", newVersion)))
}
