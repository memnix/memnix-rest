package mcq

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
)

type IUseCase interface {
	// GetByID returns the mcq with the given id.
	GetByID(ctx context.Context, id uint) (domain.Mcq, error)
}

type IRepository interface {
	// GetByID returns the mcq with the given id.
	GetByID(ctx context.Context, id uint) (domain.Mcq, error)
}

type IRedisRepository interface {
	// GetByID returns the mcq with the given id.
	GetByID(ctx context.Context, id uint) (string, error)
}
