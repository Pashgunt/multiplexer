package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
	"transport/internal/application/observability/logging"
	"transport/internal/domain/app"
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/infrastructure/config"
	"transport/internal/messaging/kafka"

	kafkago "github.com/segmentio/kafka-go"
)

func main() {
	ctxGracefulShutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defaultApp := app.DefaultApp()
	environment := &config.Environment{}
	_ = environment.Init()
	cfg, _ := config.NewLoader(&config.Validator{}, environment).Load("./configs/transport.yaml")
	adapter := kafka.NewAdapter(*cfg, logging.NewKafkaConnectionLogger(slog.LevelDebug))
	adapter.ConnectAll(kafkaconnection.DefaultKafkaConn())

	if err := adapter.WaitForConnections(defaultApp.TimeoutConnection()); err != nil {
		fmt.Println(err.Error())
	}

	for _, connection := range adapter.Connections() {
		go func() {
			connection.StartConsumers()

			for {
				message := connection.Consumer().Fetch()
				connection.Consumer().Commit([]kafkago.Message{message}, kafkaconnection.DefaultConsumer())
			}
		}()
	}

	<-ctxGracefulShutdown.Done()
	ctxShutdown, stopCtxShutdown := context.WithTimeout(context.Background(), 1*time.Second)
	adapter.CloseAll()
	defer func() {
		stop()
		stopCtxShutdown()
	}()
	<-ctxShutdown.Done()
}
