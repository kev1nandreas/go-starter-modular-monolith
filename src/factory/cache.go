package factory

import (
	"github.com/go-redis/redis/v8"
	"github.com/kev1nandreas/go-rest-api-template/config"
)

func NewRedisClient() *redis.Client {
	redisHost := config.GetEnvString("REDIS_HOST", "localhost")
	redisPort := config.GetEnvString("REDIS_PORT", "6379")
	redisPassword := config.GetEnvString("REDIS_PASSWORD", "")

	return redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})
}
