package card

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
)

type UseCase struct {
	IRepository
	IRedisRepository
}

func NewUseCase(repo IRepository, redis IRedisRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis}
}

func (*UseCase) GetByID(_ context.Context, _ uint) (domain.Card, error) {
	return domain.Card{}, nil
}
