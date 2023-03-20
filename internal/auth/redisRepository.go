package auth

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisRepository is the redis repository for the auth use case.
type RedisRepository struct {
	RedisConn *redis.Client
}

// NewRedisRepository returns a new redis repository.
func NewRedisRepository(redisConn *redis.Client) IAuthRedisRepository {
	return &RedisRepository{
		RedisConn: redisConn,
	}
}

// HasState checks if the state exists in the redis database.
func (r *RedisRepository) HasState(state string) (bool, error) {
	return r.RedisConn.SIsMember(context.Background(), "state", state).Result()
}

// SetState sets the state in the redis database.
func (r *RedisRepository) SetState(state string) error {
	return r.RedisConn.SAdd(context.Background(), "state", state).Err()
}

// DeleteState deletes the state in the redis database.
func (r *RedisRepository) DeleteState(state string) error {
	return r.RedisConn.SRem(context.Background(), "state", state).Err()
}
