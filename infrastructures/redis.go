package infrastructures

import (
	"context"

	"github.com/memnix/memnix-rest/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var redisClient *redis.Client

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
	var redisHost string
	if config.IsDevelopment() {
		redisHost = config.EnvHelper.GetEnv("DEBUG_REDIS_URL")
	} else {
		redisHost = config.EnvHelper.GetEnv("REDIS_URL")
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: config.RedisMinIdleConns,
		PoolSize:     config.RedisPoolSize,
		PoolTimeout:  config.RedisPoolTimeout,
	})

	if err := redisotel.InstrumentTracing(client); err != nil {
		log.Error().Err(err).Msg("Failed to instrument redis")
	}

	return client
}
