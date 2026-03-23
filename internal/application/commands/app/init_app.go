package appcommand

import (
	"context"
	"errors"
	"net/http"
	"transport/api/src/public"
	appconfig "transport/internal/infrastructure/config/app"
	"transport/internal/infrastructure/db"
	appredis "transport/internal/infrastructure/redis"
	"transport/pkg/logging"
	"transport/pkg/utils/backoff"
)

type IApp interface {
	StartAll(ctx context.Context)
	StopAll(ctx context.Context)
}

type App struct {
	http   public.IHttpServer
	logger logging.AdapterInterface
	redis  appredis.IRedis
	pgsql  db.IDB
}

func (a App) HTTP() public.IHttpServer {
	return a.http
}

func (a App) Logger() logging.AdapterInterface {
	return a.logger
}

func (a App) Redis() appredis.IRedis {
	return a.redis
}

func (a App) Pgsql() db.IDB {
	return a.pgsql
}

func NewApp(config appconfig.Config) App {
	logger := logging.NewAdapter(map[backoff.LoggerType]backoff.LoggerLevel{
		backoff.KafkaLogger: backoff.LoggerLevel(config.Environment.Get(backoff.EnvKafkaDebugLevelKey)),
		backoff.AppLogger:   backoff.LoggerLevel(config.Environment.Get(backoff.EnvAppDebugLevelKey)),
		backoff.APILogger:   backoff.LoggerLevel(config.Environment.Get(backoff.EnvAPIDebugLevelKey)),
	})
	logger.Init([]backoff.LoggerType{backoff.KafkaLogger, backoff.AppLogger, backoff.APILogger})

	return App{
		http: public.NewHTTPServer(
			config,
			logger.GetLogger(backoff.APILogger),
		),
		logger: logger,
		redis:  appredis.NewRedis(appredis.NewParams(config.Environment)),
		pgsql:  db.NewPostgresSQLDB(db.NewParams(config.Environment)),
	}
}

func (a App) StartAll(_ context.Context) {
	if err := a.pgsql.Open(); err != nil {
		panic(err)
	}

	if err := a.redis.Ping(); err != nil {
		panic(err)
	}

	a.http.HandleFunc(a.pgsql)

	go func() {
		if err := a.http.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()
}

func (a App) StopAll(ctx context.Context) {
	if err := a.redis.Close(); err != nil {
		panic(err)
	}

	if err := a.http.Shutdown(ctx); err != nil {
		panic(err)
	}
}
