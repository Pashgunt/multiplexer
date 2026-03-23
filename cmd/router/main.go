package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
	appcommand "transport/internal/application/commands/app"
)

func main() {
	ctxGracefulShutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := appcommand.NewApp(appcommand.
		NewKernel().
		Init().
		Config())

	app.StartAll(ctxGracefulShutdown)

	<-ctxGracefulShutdown.Done()
	ctxShutdown, stopCtxShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCtxShutdown()

	app.StopAll(ctxShutdown)
}
