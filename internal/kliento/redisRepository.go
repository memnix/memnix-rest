package kliento

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	RedisConn *redis.Client
}

func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return &RedisRepository{RedisConn: redisConn}
}

// GetName returns the name of the repository.
func (r *RedisRepository) GetName() string {
	return r.RedisConn.Get(context.Background(), "name").Val()
}

// SetName sets the name of the repository.
func (r *RedisRepository) SetName(name string) error {
	return r.RedisConn.Set(context.Background(), "name", name, 0).Err()
}
