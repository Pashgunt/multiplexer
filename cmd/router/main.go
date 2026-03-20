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
)

func main() {
	//todo добавить сначал проверкич то все запустилось и только потом чтобы начинало все работать
	ctxGracefulShutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	app := appcommand.NewKernel().Init().Config()
	adapter := kafka.NewAdapter(app)
	adapter.ConnectAll(kafkaconnection.DefaultKafkaConn())

	if !adapter.HasConnection() {
		stop()
	}

	go kafkacommand.StartProcess(adapter.Connections(), app)

	httpServer := public.NewHttpServer(app)

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
