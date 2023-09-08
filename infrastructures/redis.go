package infrastructures

import (
	"context"
	"time"

	"github.com/memnix/memnix-rest/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

var redisClient *redis.Client

// ConnectRedis Connects to redis
func ConnectRedis(redisConf config.RedisConfigStruct) error {
	redisClient = NewRedisClient(redisConf)

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
func NewRedisClient(redisConf config.RedisConfigStruct) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		MinIdleConns: redisConf.MinIdleConns,
		PoolSize:     redisConf.PoolSize,
		PoolTimeout:  time.Duration(redisConf.PoolTimeout) * time.Second,
	})

	if err := redisotel.InstrumentTracing(client); err != nil {
		otelzap.L().Error("failed to instrument redis", zap.Error(err))
	}

	return client
}
