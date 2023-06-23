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
func (r *RedisRepository) HasState(ctx context.Context, state string) (bool, error) {
	return r.RedisConn.SIsMember(ctx, "state", state).Result()
}

// SetState sets the state in the redis database.
func (r *RedisRepository) SetState(ctx context.Context, state string) error {
	return r.RedisConn.SAdd(ctx, "state", state).Err()
}

// DeleteState deletes the state in the redis database.
func (r *RedisRepository) DeleteState(ctx context.Context, state string) error {
	return r.RedisConn.SRem(ctx, "state", state).Err()
}
