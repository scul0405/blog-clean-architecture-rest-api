package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/scul0405/blog-clean-architecture-rest-api/config"
	"time"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.RedisHost,
		Password:     cfg.Redis.RedisPassword,
		DB:           cfg.Redis.RedisDb,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
	})

	return client
}
