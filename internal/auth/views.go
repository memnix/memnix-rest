package auth

import "github.com/memnix/memnix-rest/domain"

type LoginTokenVM struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func NewLoginTokenVM(token string, message string) LoginTokenVM {
	return LoginTokenVM{
		Token:   token,
		Message: message,
	}
}

type RegisterVM struct {
	Message string            `json:"message"`
	User    domain.PublicUser `json:"user"`
}

func NewRegisterVM(message string, user domain.PublicUser) RegisterVM {
	return RegisterVM{
		Message: message,
		User:    user,
	}
}
