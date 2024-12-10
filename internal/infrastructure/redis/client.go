package redis

import "github.com/redis/go-redis/v9"

type RedisClient struct {
	Client *redis.Client
}

func NewRedis() RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		Protocol: 0, // use default DB
	})
	return RedisClient{
		Client: rdb,
	}
}
