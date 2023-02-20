//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/memnix/memnix-rest/app/http/controllers"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/internal/kliento"
	"github.com/memnix/memnix-rest/internal/user"
)

func InitializeKliento() controllers.KlientoController {
	wire.Build(controllers.NewKlientoController, kliento.NewUseCase, kliento.NewRedisRepository, infrastructures.GetRedisClient)
	return controllers.KlientoController{}
}

func InitializeUser() controllers.UserController {
	wire.Build(controllers.NewUserController, user.NewUseCase, user.NewEdgeRepository, infrastructures.GetEdgeDBClient)
	return controllers.UserController{}
}

func InitializeAuth() controllers.AuthController {
	wire.Build(controllers.NewAuthController, auth.NewUseCase, user.NewEdgeRepository, infrastructures.GetEdgeDBClient)
	return controllers.AuthController{}
}
