package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	redisPort          = "6379"
	redisSecondaryPort = "6380"
	redisAddr          = ":" + redisPort
	redisPassword      = "123"
)

func redisOptions() *redis.Options {
	return &redis.Options{
		Addr:               redisAddr,
		Password:           redisPassword,
		DB:                 0,
		MaxRetries:         -1,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolSize:           10,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
	}
}
