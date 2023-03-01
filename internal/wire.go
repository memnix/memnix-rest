//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/memnix/memnix-rest/app/http/controllers"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/memnix/memnix-rest/pkg/cacheset"
)

func InitializeUser() controllers.UserController {
	wire.Build(controllers.NewUserController, user.NewUseCase, user.NewRepository, infrastructures.GetDBConn)
	return controllers.UserController{}
}

func InitializeAuth() controllers.AuthController {
	wire.Build(controllers.NewAuthController, auth.NewUseCase, user.NewRepository, infrastructures.GetDBConn)
	return controllers.AuthController{}
}

func InitializeJWT() controllers.JwtController {
	wire.Build(controllers.NewJwtController, user.NewUseCase, user.NewRepository, infrastructures.GetDBConn)
	return controllers.JwtController{}
}

func InitializeOAuth() controllers.OAuthController {
	wire.Build(controllers.NewOAuthController, auth.NewUseCase, user.NewRepository, infrastructures.GetDBConn, cacheset.New)
	return controllers.OAuthController{}
}
