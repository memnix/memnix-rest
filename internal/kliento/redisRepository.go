package kliento

import "github.com/go-redis/redis/v8"

type RedisRepository struct {
	RedisConn *redis.Client
}

func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return &RedisRepository{RedisConn: redisConn}
}

// GetName returns the name of the repository.
func (r *RedisRepository) GetName() string {
	return r.RedisConn.Get(r.RedisConn.Context(), "name").Val()
}

// SetName sets the name of the repository.
func (r *RedisRepository) SetName(name string) error {
	return r.RedisConn.Set(r.RedisConn.Context(), "name", name, 0).Err()
}
