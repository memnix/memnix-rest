package user

import (
	"context"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	RedisConn *redis.Client
}

func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return &RedisRepository{
		RedisConn: redisConn,
	}
}

func (r *RedisRepository) Get(id uint) (string, error) {
	return r.RedisConn.Get(context.Background(), "user:"+utils.ConvertUIntToStr(id)).Result()
}

func (r *RedisRepository) Set(id uint, value string) error {
	return r.RedisConn.Set(context.Background(), "user:"+utils.ConvertUIntToStr(id), value, config.RedisDefaultExpireTime).Err()
}
