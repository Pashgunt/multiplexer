package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transport/internal/infrastructure/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-quit
		cancel()
	}()

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	<-ctx.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}
}
