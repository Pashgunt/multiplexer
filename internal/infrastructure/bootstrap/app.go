package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"transport/api/src/domainservice"
	"transport/api/src/factory"
	"transport/api/src/handler"
	"transport/api/src/repository"
	"transport/api/src/router"
	"transport/api/src/service"
	appconfig "transport/internal/infrastructure/config/app"
	"transport/pkg/logging"
	"transport/pkg/utils/backoff"
	"transport/providers"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		getProvideOptions()...,
	)
}

func getProvideOptions() []fx.Option {
	return []fx.Option{
		fx.Provide(
			providers.Environment,
			providers.Config,
			providers.Logger,
			providers.DB,
			providers.Cache,
			providers.HTTP,

			repository.NewTargetServiceRepository,

			factory.NewTargetServiceFactory,

			domainservice.NewTargetDomainService,

			service.NewTargetServiceService,

			handler.NewTargetServiceHandler,

			router.NewRouter,
		),
		fx.Invoke(
			runMigrations,
			startServer,
		),
	}
}

func startServer(srv *http.Server, cfg appconfig.Config, logger logging.AdapterInterface, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.GetLogger(backoff.AppLogger).Info(logging.AppLogEntity{Message: fmt.Sprintf("Starting HTTP server. Addr %s. Host: %s. Port: %s.", srv.Addr, cfg.HTTP.Host, cfg.HTTP.Port)})

				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.GetLogger(backoff.AppLogger).Error(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.GetLogger(backoff.AppLogger).Info(logging.AppLogEntity{Message: "Shutting down server..."})
			return srv.Shutdown(ctx)
		},
	})
}

func runMigrations(cfg appconfig.Config, logger logging.AdapterInterface) error {
	logger.GetLogger(backoff.AppLogger).Info(logging.AppLogEntity{Message: "Running database migrations"})

	db, err := sql.Open("postgres", cfg.DB.DatabaseSourceName)

	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			logger.GetLogger(backoff.AppLogger).Error(err)
		}
	}(db)

	if err = goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err = goose.Up(db, "migrations"); err != nil {
		return err
	}

	logger.GetLogger(backoff.AppLogger).Info(logging.AppLogEntity{Message: "Database migrations completed successfully"})

	return nil
}
