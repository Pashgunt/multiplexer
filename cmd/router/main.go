package main

import (
	"transport/internal/infrastructure/config"
)

func main() {
	environment := &config.Environment{}
	environment.Init()
	loader := config.NewLoader(&config.Validator{}, environment)
	loader.Load("./configs/transport.yaml")
}
