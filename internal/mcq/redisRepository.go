package mcq

import (
	"context"

	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func getBaseKey() string {
	return "mcq:id:"
}

func withID(getter func() string, id uint) string {
	return getter() + utils.ConvertUIntToStr(id)
}

type RedisRepository struct {
	RedisConn *redis.Client // RedisConn is the redis connection.
}

// GetByID gets the mcq by id.
func (r RedisRepository) GetByID(ctx context.Context, id uint) (string, error) {
	return r.RedisConn.Get(ctx, withID(getBaseKey, id)).Result()
}

func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return RedisRepository{RedisConn: redisConn}
}
