package db

import (
	"github.com/go-redis/redis/v8"
)

func NewRedis(url string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	if _, err := client.Ping(client.Context()).Result(); err != nil {
		panic(err)
	}

	return client
}
