package card

import "github.com/memnix/memnix-rest/domain"

type UseCase struct {
	IRepository
	IRedisRepository
}

func NewUseCase(repo IRepository, redis IRedisRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis}
}

func (u *UseCase) GetByID(id uint) (domain.Card, error) {
	return domain.Card{}, nil
}
