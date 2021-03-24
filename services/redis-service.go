package services

import (
    "context"
    "github.com/go-redis/redis/v8"
    "os"
)

var ctx = context.Background()
var rdb *redis.Client

func GetRedisService() *redis.Client{
    if rdb != nil {
        return rdb
    }

    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("RedisHostPort"),
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    return rdb
}
