package redisdb

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func CreateRedisClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       dbNo,
	})

	return rdb
}
