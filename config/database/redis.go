package database

import (
	"context"

	"github.com/GSVillas/movie-pass-api/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Env.RedisAddress,
		Password: config.Env.RedisPassword,
		DB:       config.Env.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
