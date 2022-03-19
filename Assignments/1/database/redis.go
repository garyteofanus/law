package database

import "github.com/go-redis/redis/v8"

func NewRedis(host string, port string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0, // use default DB
	})
}
