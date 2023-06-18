package card

import "github.com/memnix/memnix-rest/domain"

type IUseCase interface {
	GetByID(id uint) (domain.Card, error)
}

type IRepository interface {
	GetByID(id uint) (domain.Card, error)
}

type IRedisRepository interface {
	GetByID(id uint) (string, error)
}
