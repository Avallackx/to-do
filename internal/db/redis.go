package db

import (
	"todo-app/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func InitializeRedisConn() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisHost(),
		DB:   config.RedisDB(),
	})
}
