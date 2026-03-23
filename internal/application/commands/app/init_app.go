package appcommand

import (
	"context"
	"sync"
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
	config appconfig.Config
	logger logging.AdapterInterface
	redis  appredis.IRedis
	pgsql  db.IDB
}

func (a App) HTTP() public.IHttpServer {
	return a.http
}

func (a App) Config() appconfig.Config {
	return a.config
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
	sqldb := db.NewPostgresSQLDB(config.Environment.Get(envNamePgDatabaseURL))

	return App{
		http: public.NewHTTPServer(
			config,
			logger.GetLogger(backoff.APILogger),
			sqldb,
		),
		config: config,
		logger: logger,
		redis: appredis.NewRedis(
			config.Environment.Get("REDIS"),
			config.Environment.Get("REDIS_PASSWORD"),
		),
		pgsql: sqldb,
	}
}

func (a App) StartAll(_ context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(3)

	wg.Go(func() {
		if err := a.pgsql.Open(); err != nil {
			panic(err)
		}
	})

	wg.Go(func() {
		if err := a.redis.Ping(); err != nil {
			panic(err)
		}
	})

	wg.Go(func() {
		if err := a.http.Start(); err != nil {
			panic(err)
		}
	})

	wg.Wait()
}

func (a App) StopAll(ctx context.Context) {
	if err := a.redis.Close(); err != nil {
		panic(err)
	}

	if err := a.http.Shutdown(ctx); err != nil {
		panic(err)
	}
}
