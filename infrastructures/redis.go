package infrastructures

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/memnix/memnix-rest/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	client *redis.Client
}

var (
	redisInstance *RedisManager //nolint:gochecknoglobals //Singleton
	redisOnce     sync.Once     //nolint:gochecknoglobals //Singleton
)

func GetRedisClient() *redis.Client {
	return GetRedisManagerInstance().GetRedisClient()
}

func GetRedisManagerInstance() *RedisManager {
	redisOnce.Do(func() {
		redisInstance = &RedisManager{}
	})
	return redisInstance
}

func (r *RedisManager) ConnectRedis(redisConf config.RedisConfig) error {
	r.client = r.NewRedisClient(redisConf)

	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisManager) CloseRedis() error {
	return r.client.Close()
}

func (r *RedisManager) GetRedisClient() *redis.Client {
	return r.client
}

func (r *RedisManager) NewRedisClient(redisConf config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		MinIdleConns: redisConf.MinIdleConns,
		PoolSize:     redisConf.PoolSize,
		PoolTimeout:  time.Duration(redisConf.PoolTimeout) * time.Second,
	})

	if err := redisotel.InstrumentTracing(client); err != nil {
		log.Error("failed to instrument redis", slog.Any("error", err))
	}

	return client
}
