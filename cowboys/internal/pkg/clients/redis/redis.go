package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	rdb *redis.Client
}

func (cli *RedisClient) Set(key string, data interface{}) error {
	err := cli.rdb.Set(ctx, key, data, 0).Err()

	return err
}

func (cli *RedisClient) Get(key string) (string, error) {
	val, err := cli.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	} 

	return val, nil
}

func New(redisHost string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", redisHost),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisClient{rdb: rdb}
}
