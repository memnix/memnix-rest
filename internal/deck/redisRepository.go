package deck

import (
	"context"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/redis/go-redis/v9"
)

func getBaseKey() string {
	return "deck:id:"
}

func getOwnedKey() string {
	return "deck:owned:"
}

func getLearningKey() string {
	return "deck:learning:"
}

func withID(getter func() string, id uint) string {
	return getter() + utils.ConvertUIntToStr(id)
}

// RedisRepository is the interface for the redis repository.
type RedisRepository struct {
	RedisConn *redis.Client // RedisConn is the redis connection.
}

// GetByID gets the deck by id.
func (r RedisRepository) GetByID(id uint) (string, error) {
	return r.RedisConn.Get(context.Background(), withID(getBaseKey, id)).Result()
}

// SetByID sets the deck by id.
func (r RedisRepository) SetByID(id uint, deck string) error {
	return r.RedisConn.Set(context.Background(), withID(getBaseKey, id), deck, config.RedisDefaultExpireTime).Err()
}

// DeleteByID deletes the deck by id.
func (r RedisRepository) DeleteByID(id uint) error {
	return r.RedisConn.Del(context.Background(), withID(getBaseKey, id)).Err()
}

// SetOwnedByUser sets the decks owned by the user.
func (r RedisRepository) SetOwnedByUser(userID uint, decks string) error {
	return r.RedisConn.Set(context.Background(), withID(getOwnedKey, userID), decks, config.RedisOwnedExpireTime).Err()
}

// GetOwnedByUser gets the decks owned by the user.
func (r RedisRepository) GetOwnedByUser(userID uint) (string, error) {
	return r.RedisConn.Get(context.Background(), withID(getOwnedKey, userID)).Result()
}

// DeleteOwnedByUser deletes the decks owned by the user.
func (r RedisRepository) DeleteOwnedByUser(userID uint) error {
	return r.RedisConn.Del(context.Background(), withID(getOwnedKey, userID)).Err()
}

// SetLearningByUser sets the decks learning by the user.
func (r RedisRepository) SetLearningByUser(userID uint, decks string) error {
	return r.RedisConn.Set(context.Background(), withID(getLearningKey, userID), decks, config.RedisOwnedExpireTime).Err()
}

// GetLearningByUser gets the decks learning by the user.
func (r RedisRepository) GetLearningByUser(userID uint) (string, error) {
	return r.RedisConn.Get(context.Background(), withID(getLearningKey, userID)).Result()
}

// DeleteLearningByUser deletes the decks learning by the user.
func (r RedisRepository) DeleteLearningByUser(userID uint) error {
	return r.RedisConn.Del(context.Background(), withID(getLearningKey, userID)).Err()
}

// NewRedisRepository returns a new redis repository.
func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return RedisRepository{
		RedisConn: redisConn,
	}
}
