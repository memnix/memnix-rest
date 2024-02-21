package deck

import (
	"context"
	"time"

	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/redis/go-redis/v9"
)

const (
	OwnedExpirationTime = 2 * time.Hour
	defaultExpireTime   = 6 * time.Hour
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
func (r RedisRepository) GetByID(ctx context.Context, id uint) (string, error) {
	return r.RedisConn.Get(ctx, withID(getBaseKey, id)).Result()
}

// SetByID sets the deck by id.
func (r RedisRepository) SetByID(ctx context.Context, id uint, deck string) error {
	return r.RedisConn.Set(ctx, withID(getBaseKey, id), deck, defaultExpireTime).Err()
}

// DeleteByID deletes the deck by id.
func (r RedisRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.RedisConn.Del(ctx, withID(getBaseKey, id)).Err()
}

// SetOwnedByUser sets the decks owned by the user.
func (r RedisRepository) SetOwnedByUser(ctx context.Context, userID uint, decks string) error {
	return r.RedisConn.Set(ctx, withID(getOwnedKey, userID), decks, OwnedExpirationTime).Err()
}

// GetOwnedByUser gets the decks owned by the user.
func (r RedisRepository) GetOwnedByUser(ctx context.Context, userID uint) (string, error) {
	return r.RedisConn.Get(ctx, withID(getOwnedKey, userID)).Result()
}

// DeleteOwnedByUser deletes the decks owned by the user.
func (r RedisRepository) DeleteOwnedByUser(ctx context.Context, userID uint) error {
	return r.RedisConn.Del(ctx, withID(getOwnedKey, userID)).Err()
}

// SetLearningByUser sets the decks learning by the user.
func (r RedisRepository) SetLearningByUser(ctx context.Context, userID uint, decks string) error {
	return r.RedisConn.Set(ctx, withID(getLearningKey, userID), decks, OwnedExpirationTime).Err()
}

// GetLearningByUser gets the decks learning by the user.
func (r RedisRepository) GetLearningByUser(ctx context.Context, userID uint) (string, error) {
	return r.RedisConn.Get(ctx, withID(getLearningKey, userID)).Result()
}

// DeleteLearningByUser deletes the decks learning by the user.
func (r RedisRepository) DeleteLearningByUser(ctx context.Context, userID uint) error {
	return r.RedisConn.Del(ctx, withID(getLearningKey, userID)).Err()
}

// NewRedisRepository returns a new redis repository.
func NewRedisRepository(redisConn *redis.Client) IRedisRepository {
	return RedisRepository{
		RedisConn: redisConn,
	}
}
