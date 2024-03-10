package user

import (
	"context"
	"time"

	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/redis/go-redis/v9"
)

const defaultExpireTime = 6 * time.Hour

// RedisRepository is the interface for the redis repository.
type RedisRepository struct {
	RedisConn *redis.Client // RedisConn is the redis connection.
}

// NewRedisRepository returns a new redis repository.
func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return &RedisRepository{
		RedisConn: redisConn,
	}
}

// Get gets the user by id.
func (r *RedisRepository) Get(ctx context.Context, id int32) (string, error) {
	return r.RedisConn.Get(ctx, keyPrefix+utils.ConvertInt32ToStr(id)).Result()
}

// Set sets the user by id.
func (r *RedisRepository) Set(ctx context.Context, id int32, value string) error {
	return r.RedisConn.Set(ctx, keyPrefix+utils.ConvertInt32ToStr(id), value, defaultExpireTime).Err()
}
