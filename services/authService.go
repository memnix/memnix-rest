package services

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"github.com/memnix/memnixrest/interfaces"
)

type AuthService struct {
	interfaces.IAuthService
}

func (k *kernel) InjectAuthController() controllers.AuthController {
	DBConn := infrastructures.GetDBConn()

	authRepository := &repositories.AuthRepository{DBConn: DBConn}
	authService := &AuthService{IAuthService: authRepository}
	authController := controllers.AuthController{IAuthService: authService}

	return authController
}
