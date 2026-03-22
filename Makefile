start-app:
	go run transport/cmd/router

migrate-up:
	go run transport/cmd/migrate

migrate-status:
	go run transport/cmd/migrate --status

migrate-down:
	go run transport/cmd/migrate --down

migrate-reset:
	go run transport/cmd/migrate --reset

linter:
	golangci-lint run --fix