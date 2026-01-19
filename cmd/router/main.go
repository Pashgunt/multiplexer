package main

import "transport/internal/infrastructure/config"

func main() {
	environment := &config.Environment{}
	environment.Init()
	loader := config.Loader{Validator: &config.Validator{}, Environment: environment}
	loader.Load("./configs/transport.yaml")
}
