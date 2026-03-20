package appcommand

import (
	"context"
	"transport/api/src/public"
	kafkacommand "transport/internal/application/commands/kafka"
	kafkaconnection "transport/internal/domain/connection"
	appconfig "transport/internal/infrastructure/config/app"
	"transport/internal/messaging/kafka"
	"transport/pkg/logging"
	"transport/pkg/utils/backoff"
)

type IApp interface {
	StartAll(ctx context.Context)
	StopAll(ctx context.Context)
}

type App struct {
	http   public.IHttpServer
	kafka  kafka.AdapterInterface
	config appconfig.Config
	logger logging.LoggerInterface
}

func NewApp(config appconfig.Config) App {
	appLogger := config.Logger.GetLogger(backoff.AppLogger)
	appLogger.Info(logging.NewAppLogEntity("kernel config initialized"))

	appLogger.Info(logging.NewAppLogEntity("Start load Kafka connections"))
	adapter := kafka.NewAdapter(config)
	adapter.ConnectAll(kafkaconnection.DefaultKafkaConn())
	appLogger.Info(logging.NewAppLogEntity("Loaded Kafka connections"))

	appLogger.Info(logging.NewAppLogEntity("Start Http Server"))
	httpServer := public.NewHttpServer(config)

	appLogger.Info(logging.NewAppLogEntity("Started Http Server"))

	return App{
		http:   httpServer,
		kafka:  adapter,
		config: config,
		logger: appLogger,
	}
}

func (a App) StartAll(ctx context.Context) {
	go kafkacommand.StartProcess(a.kafka.Connections(), a.config)
	go a.http.Start()
}

func (a App) StopAll(ctx context.Context) {
	a.logger.Info(logging.NewAppLogEntity("start close all connections"))
	a.kafka.CloseAll(ctx)
	a.http.Shutdown(ctx)
	a.logger.Info(logging.NewAppLogEntity("shutdown http server connection"))
}
