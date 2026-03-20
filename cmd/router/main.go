package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
	"transport/api/src/public"
	appcommand "transport/internal/application/commands/app"
	kafkacommand "transport/internal/application/commands/kafka"
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/messaging/kafka"
	"transport/pkg/logging"
	"transport/pkg/utils/backoff"
)

func main() {
	ctxGracefulShutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	app := appcommand.NewKernel().Init().Config()

	appLogger := app.Logger.GetLogger(backoff.AppLogger)
	appLogger.Info(logging.NewAppLogEntity("kernel config initialized"))

	appLogger.Info(logging.NewAppLogEntity("Start load Kafka connections"))
	adapter := kafka.NewAdapter(app)
	adapter.ConnectAll(kafkaconnection.DefaultKafkaConn())

	if !adapter.HasConnection() {
		stop()
	}

	appLogger.Info(logging.NewAppLogEntity("Loaded Kafka connections"))

	appLogger.Info(logging.NewAppLogEntity("Start Http Server"))
	httpServer := public.NewHttpServer(app)

	appLogger.Info(logging.NewAppLogEntity("Started Http Server"))

	go kafkacommand.StartProcess(adapter.Connections(), app)
	go httpServer.Start()

	<-ctxGracefulShutdown.Done()
	ctxShutdown, stopCtxShutdown := context.WithTimeout(context.Background(), 10*time.Second)

	adapter.CloseAll(ctxShutdown)
	httpServer.Shutdown(ctxShutdown)

	defer func() {
		stop()
		stopCtxShutdown()
	}()
}
