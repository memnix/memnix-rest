package mcq

import "github.com/memnix/memnix-rest/domain"

type IUseCase interface {
	// GetByID returns the mcq with the given id.
	GetByID(id uint) (domain.Mcq, error)
}

type IRepository interface {
	// GetByID returns the mcq with the given id.
	GetByID(id uint) (domain.Mcq, error)
}

type IRedisRepository interface {
	// GetByID returns the mcq with the given id.
	GetByID(id uint) (string, error)
}
