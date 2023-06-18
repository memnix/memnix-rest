package mcq

import "github.com/memnix/memnix-rest/domain"

type UseCase struct {
	IRepository
	IRedisRepository
}

func (c *UseCase) GetByID(id uint) (domain.Mcq, error) {
	return c.IRepository.GetByID(id)
}

func NewUseCase(repo IRepository, redis IRedisRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis}
}
