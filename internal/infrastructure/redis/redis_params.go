package redis

import "transport/internal/infrastructure/config"

const (
	addr     = "REDIS"
	password = "REDIS_PASSWORD"
)

type Params struct {
	Addr     string
	Password string
	DB       int
}

func NewParams(env config.EnvironmentInterface) Params {
	return Params{
		Addr:     env.Get(addr),
		Password: env.Get(password),
	}
}
