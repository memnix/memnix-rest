package user

import (
	"context"

	"github.com/fxamacker/cbor/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/pkg/utils"
)

// UseCase is the user use case.
type UseCase struct {
	IRepository
	IRedisRepository
}

// GetName returns the name of the user with the given id.
func (u UseCase) GetName(ctx context.Context, id string) string {
	uintID, _ := utils.ConvertStrToUInt(id)

	return u.IRepository.GetName(ctx, uintID)
}

// GetByID returns the user with the given id.
func (u UseCase) GetByID(ctx context.Context, id uint) (domain.User, error) {
	var userObject domain.User

	if cacheHit, _ := u.IRedisRepository.Get(ctx, id); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &userObject); err == nil {
			return userObject, nil
		}
	}

	userObject, err := u.IRepository.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	if marshalledUser, err := cbor.Marshal(userObject); err == nil {
		_ = u.IRedisRepository.Set(ctx, id, string(marshalledUser))
	}
	return userObject, nil
}

// NewUseCase returns a new user use case.
func NewUseCase(repo IRepository, redis IRedisRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis}
}
