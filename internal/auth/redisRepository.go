package auth

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	RedisConn *redis.Client
}

func NewRedisRepository(redisConn *redis.Client) IAuthRedisRepository {
	return &RedisRepository{
		RedisConn: redisConn,
	}
}

func (r *RedisRepository) HasState(state string) (bool, error) {
	return r.RedisConn.SIsMember(context.Background(), "state", state).Result()
}

func (r *RedisRepository) SetState(state string) error {
	return r.RedisConn.SAdd(context.Background(), "state", state).Err()
}

func (r *RedisRepository) DeleteState(state string) error {
	return r.RedisConn.SRem(context.Background(), "state", state).Err()
}
