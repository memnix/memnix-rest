//go:build wireinject
// +build wireinject

package services

import (
	"github.com/google/wire"
	"github.com/memnix/memnix-rest/app/v2/handlers"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/services/auth"
	"github.com/memnix/memnix-rest/services/user"
)

func InitializeAuthHandler() handlers.AuthController {
	wire.Build(handlers.NewAuthController, auth.NewUseCase, user.NewRepository, infrastructures.GetPgxConn)
	return handlers.AuthController{}
}
