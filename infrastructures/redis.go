package infrastructures

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	client      *redis.Client
	Redisconfig RedisConfig
}

type RedisConfig struct {
	Addr         string
	Password     string
	MinIdleConns int
	PoolSize     int
	PoolTimeout  int
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

func NewRedisInstance(redisConf RedisConfig) *RedisManager {
	return GetRedisManagerInstance().RedisWithConfig(redisConf)
}

func (r *RedisManager) RedisWithConfig(redisConf RedisConfig) *RedisManager {
	r.Redisconfig = redisConf
	return r
}

func (r *RedisManager) ConnectRedis() error {
	r.client = redis.NewClient(&redis.Options{
		Addr:         r.Redisconfig.Addr,
		Password:     r.Redisconfig.Password,
		MinIdleConns: r.Redisconfig.MinIdleConns,
		PoolSize:     r.Redisconfig.PoolSize,
		PoolTimeout:  time.Duration(r.Redisconfig.PoolTimeout) * time.Second,
	})

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
