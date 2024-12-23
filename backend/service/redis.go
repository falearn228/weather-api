package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/weather-api/configs"
)

func NewRedisClient() *redis.Client {
	config := configs.LoadConfig()
	return redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		Password: config.RedisPassword,
		DB: 0,
	})
}