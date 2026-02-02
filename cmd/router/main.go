package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
	kafkacommand "transport/internal/application/commands/kafka"
	"transport/internal/application/observability/logging"
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/infrastructure/config"
	"transport/internal/messaging/kafka"
)

func main() {
	ctxGracefulShutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	env, _ := config.NewEnvironment()

	cfg, _ := config.NewLoader(config.NewValidator(), env).Load("./configs/transport.yaml")
	logger := logging.NewKafkaConnectionLogger(slog.LevelDebug)
	adapter := kafka.NewAdapter(*cfg, logger)
	adapter.ConnectAll(kafkaconnection.DefaultKafkaConn())

	kafkacommand.StartProcess(adapter.Connections(), logger)

	<-ctxGracefulShutdown.Done()
	ctxShutdown, stopCtxShutdown := context.WithTimeout(context.Background(), 1*time.Second)
	adapter.CloseAll()
	defer func() {
		stop()
		stopCtxShutdown()
	}()
	<-ctxShutdown.Done()
}
