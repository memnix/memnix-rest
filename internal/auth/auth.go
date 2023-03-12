package auth

import "github.com/memnix/memnix-rest/domain"

type IUseCase interface {
	Login(password string, email string) (string, error)
	Register(registerStruct domain.Register) (domain.User, error)
	Logout() (string, error)
	RefreshToken(user domain.User) (string, error)
	RegisterOauth(user domain.User) error
	LoginOauth(user domain.User) (string, error)
}

type IAuthRedisRepository interface {
	HasState(state string) (bool, error)
	SetState(state string) error
	DeleteState(state string) error
}
