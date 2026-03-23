package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	Ping() error
	Close() error
}

type Redis struct {
	client *redis.Client
	params Params
}

func NewRedis(params Params) IRedis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     params.Addr,
			Password: params.Password,
			DB:       params.DB,
		}),
		params: params,
	}
}

func (r *Redis) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.client.Ping(ctx).Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
