package auth

import "github.com/memnix/memnix-rest/domain"

type IUseCase interface {
	Login(password string, email string) (string, error)
	Register(user domain.User) (domain.User, error)
	Logout() (string, error)
	RefreshToken(user domain.User) (string, error)
}
