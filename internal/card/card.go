package card

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
)

type IUseCase interface {
	GetByID(ctx context.Context, id uint) (domain.Card, error)
}

type IRepository interface {
	GetByID(ctx context.Context, id uint) (domain.Card, error)
}

type IRedisRepository interface {
	GetByID(ctx context.Context, id uint) (string, error)
}
