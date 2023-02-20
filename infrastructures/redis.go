package infrastructures

import (
	"github.com/memnix/memnix-rest/config"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var (
	redisClient *redis.Client
)

// ConnectRedis Connects to redis
func ConnectRedis() error {
	redisClient = NewRedisClient()

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

// CloseRedis Closes redis connection
func CloseRedis() error {
	return redisClient.Close()
}

// GetRedisClient Returns redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// NewRedisClient Returns new redis client
func NewRedisClient() *redis.Client {
	redisHost := config.RedisHost

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: config.RedisMinIdleConns,
		PoolSize:     config.RedisPoolSize,
		PoolTimeout:  config.RedisPoolTimeout,
	})

	return client
}
