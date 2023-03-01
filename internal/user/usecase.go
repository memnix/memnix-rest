package user

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/pkg/utils"
)

type UseCase struct {
	IRepository
	IRedisRepository
}

func (u UseCase) GetName(id string) string {
	uintID, _ := utils.ConvertStrToUInt(id)

	return u.IRepository.GetName(uintID)
}

func (u UseCase) GetByID(id uint) (domain.User, error) {
	var userObject domain.User

	if cacheHit, _ := u.IRedisRepository.Get(id); cacheHit != "" {
		if err := cbor.Unmarshal([]byte(cacheHit), &userObject); err == nil {
			return userObject, nil
		}
	}

	userObject, err := u.IRepository.GetByID(id)
	if err != nil {
		return domain.User{}, err
	}
	if marshalledUser, err := cbor.Marshal(userObject); err == nil {
		_ = u.IRedisRepository.Set(id, string(marshalledUser))
	}
	return userObject, nil
}

func NewUseCase(repo IRepository, redis IRedisRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis}
}
