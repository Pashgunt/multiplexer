package main

import (
	"fmt"
	"time"
	kafkaconnection "transport/internal/domain/connection"
	"transport/internal/infrastructure/config"
	"transport/internal/messaging/kafka"
)

func main() {
	environment := &config.Environment{}
	environment.Init()
	loader := config.NewLoader(&config.Validator{}, environment)
	cfg, _ := loader.Load("./configs/transport.yaml")
	adapter := kafka.NewAdapter(*cfg)
	adapter.ConnectAll(kafkaconnection.DefaultKafkaConn())
	time.Sleep(10 * time.Second)
	conns := adapter.Connections()
	fmt.Println("conns: ", conns)
}
