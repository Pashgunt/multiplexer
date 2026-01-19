package main

import "transport/internal/infrastructure/config"

func main() {
	loader := config.Loader{Validator: &config.Validator{}}
	loader.Load("./configs/transport.yaml")
}
