package redis

import "github.com/redis/go-redis/v9"

type Redis struct {
	client *redis.Client
}

func NewRedis() Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		Protocol: 0, // use default DB
	})
	return Redis{
		client: rdb,
	}
}
