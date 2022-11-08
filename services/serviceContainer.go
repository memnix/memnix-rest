package services

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"sync"
)

type IServiceContainer interface {
	InjectUserController() controllers.UserController
}

type kernel struct{}

func (k *kernel) InjectUserController() controllers.UserController {
	DBConn := infrastructures.GetDBConn()

	userRepository := &repositories.UserRepository{DBConn: DBConn}
	userService := &UserService{IUserRepository: userRepository}
	userController := controllers.UserController{IUserService: userService}

	return userController
}

var (
	k             *kernel
	containerOnce sync.Once
)

func GetServiceContainer() IServiceContainer {
	containerOnce.Do(func() {
		k = &kernel{}
	})
	return k
}
