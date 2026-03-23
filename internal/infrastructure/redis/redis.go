package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	Ping() error
	Close() error
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration)
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

func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Bytes()

	if err != nil {
		return err
	}

	return json.Unmarshal(val, dest)
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	jsonData, err := json.Marshal(value)

	if err != nil {
		return
	}

	r.client.Set(ctx, key, jsonData, expiration)
}
