package infrastructures_test

import (
	"testing"

	"github.com/memnix/memnix-rest/infrastructures"
)

func TestGetRedisClient(t *testing.T) {
	// Initialize RedisManager instance
	redisConf := infrastructures.RedisConfig{
		Addr:         "localhost:6379",
		Password:     "",
		MinIdleConns: 10,
		PoolSize:     100,
		PoolTimeout:  30,
	}
	redisManager := infrastructures.NewRedisInstance(redisConf)

	// Connect to Redis
	err := redisManager.ConnectRedis()
	if err != nil {
		t.Errorf("Failed to connect to Redis: %v", err)
	}

	// Get Redis client
	client := redisManager.GetRedisClient()
	if client == nil {
		t.Error("Failed to get Redis client")
	}

	// Close Redis connection
	err = redisManager.CloseRedis()
	if err != nil {
		t.Errorf("Failed to close Redis connection: %v", err)
	}
}

func TestRedisManager_RedisWithConfig(t *testing.T) {
	// Initialize RedisManager instance
	redisManager := &infrastructures.RedisManager{}

	// Set Redis config
	redisConf := infrastructures.RedisConfig{
		Addr:         "localhost:6379",
		Password:     "",
		MinIdleConns: 10,
		PoolSize:     100,
		PoolTimeout:  30,
	}
	redisManager.RedisWithConfig(redisConf)

	// Check if Redis config is set correctly
	if redisManager.Redisconfig != redisConf {
		t.Errorf("Redis config not set correctly")
	}
}
