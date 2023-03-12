package deck

import (
	"context"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	RedisConn *redis.Client
}

func (r RedisRepository) GetByID(id uint) (string, error) {
	return r.RedisConn.Get(context.Background(), "deck:id:"+utils.ConvertUIntToStr(id)).Result()
}

func (r RedisRepository) SetByID(id uint, deck string) error {
	return r.RedisConn.Set(context.Background(), "deck:id:"+utils.ConvertUIntToStr(id), deck, config.RedisDefaultExpireTime).Err()
}

func (r RedisRepository) DeleteByID(id uint) error {
	return r.RedisConn.Del(context.Background(), "deck:id:"+utils.ConvertUIntToStr(id)).Err()
}

func (r RedisRepository) SetOwnedByUser(userID uint, decks string) error {
	return r.RedisConn.Set(context.Background(), "deck:owned:"+utils.ConvertUIntToStr(userID), decks, config.RedisOwnedExpireTime).Err()
}

func (r RedisRepository) GetOwnedByUser(userID uint) (string, error) {
	return r.RedisConn.Get(context.Background(), "deck:owned:"+utils.ConvertUIntToStr(userID)).Result()
}

func (r RedisRepository) DeleteOwnedByUser(userID uint) error {
	return r.RedisConn.Del(context.Background(), "deck:owned:"+utils.ConvertUIntToStr(userID)).Err()
}

func (r RedisRepository) SetLearningByUser(userID uint, decks string) error {
	return r.RedisConn.Set(context.Background(), "deck:learning:"+utils.ConvertUIntToStr(userID), decks, config.RedisOwnedExpireTime).Err()
}

func (r RedisRepository) GetLearningByUser(userID uint) (string, error) {
	return r.RedisConn.Get(context.Background(), "deck:learning:"+utils.ConvertUIntToStr(userID)).Result()
}

func (r RedisRepository) DeleteLearningByUser(userID uint) error {
	return r.RedisConn.Del(context.Background(), "deck:learning:"+utils.ConvertUIntToStr(userID)).Err()
}

func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return RedisRepository{
		RedisConn: redisConn,
	}
}
