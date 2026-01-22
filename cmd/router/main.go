package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctxGracefulShutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	//services

	<-ctxGracefulShutdown.Done()
	ctxShutdown, stopCtxShutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		stop()
		stopCtxShutdown()
	}()
	<-ctxShutdown.Done()
}
