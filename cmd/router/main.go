package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
	appcommand "transport/internal/application/commands/app"
	kafkacommand "transport/internal/application/commands/kafka"
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/messaging/kafka"
)

func main() {
	ctxGracefulShutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	app := appcommand.NewKernel().Init().Config()
	adapter := kafka.NewAdapter(app)
	adapter.ConnectAll(kafkaconnection.DefaultKafkaConn())

	if adapter.HasConnection() {
		kafkacommand.StartProcess(adapter.Connections(), app)
	}

	<-ctxGracefulShutdown.Done()
	ctxShutdown, stopCtxShutdown := context.WithTimeout(context.Background(), 1*time.Second)
	adapter.CloseAll(ctxShutdown)
	defer func() {
		stop()
		stopCtxShutdown()
	}()
}
