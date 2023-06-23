package mcq

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
)

type UseCase struct {
	IRepository
	IRedisRepository
}

func (c *UseCase) GetByID(ctx context.Context, id uint) (domain.Mcq, error) {
	return c.IRepository.GetByID(ctx, id)
}

func NewUseCase(repo IRepository, redis IRedisRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis}
}
