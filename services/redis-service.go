package services

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var ctx = context.Background()
var redisClient *redis.Client

func GetRedisService() *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	_redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisHostPort"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisClient = _redisClient

	return redisClient
}
