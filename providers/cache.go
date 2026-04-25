package providers

import (
	appconfig "transport/internal/infrastructure/config/app"

	"github.com/redis/go-redis/v9"
)

func Cache(cfg appconfig.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}
