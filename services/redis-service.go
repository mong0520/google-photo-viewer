package services

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"os"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedisService(redisHostPort string) error {
	_redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisHostPort"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if _redisClient != nil {
		redisClient = _redisClient
		return nil
	} else {
		return errors.New("unable to initialize redis service")
	}
}

func GetRedisService() *redis.Client {
	return redisClient
}
