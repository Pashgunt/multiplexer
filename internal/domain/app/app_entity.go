package app

import "time"

type App struct {
	timeoutConnection time.Duration
}

func (app App) TimeoutConnection() time.Duration {
	return app.timeoutConnection
}

func DefaultApp() App {
	return App{
		timeoutConnection: 30 * time.Second,
	}
}
