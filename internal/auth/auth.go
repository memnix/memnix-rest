package auth

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
)

// IUseCase is the interface for the auth use case.
type IUseCase interface {
	Login(ctx context.Context, password string, email string) (string, error)
	Register(ctx context.Context, registerStruct domain.Register) (domain.User, error)
	Logout(ctx context.Context) (string, error)
	RefreshToken(ctx context.Context, user domain.User) (string, error)
	RegisterOauth(ctx context.Context, user domain.User) error
	LoginOauth(ctx context.Context, user domain.User) (string, error)
}

// IAuthRedisRepository is the interface for the auth redis repository.
type IAuthRedisRepository interface {
	HasState(ctx context.Context, state string) (bool, error)
	SetState(ctx context.Context, state string) error
	DeleteState(ctx context.Context, state string) error
}
