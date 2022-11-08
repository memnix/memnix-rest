package services

import (
	"errors"
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"github.com/memnix/memnixrest/interfaces"
	"github.com/memnix/memnixrest/models"
)

type UserService struct {
	interfaces.IUserRepository
}

func (k *kernel) InjectUserController() controllers.UserController {
	DBConn := infrastructures.GetDBConn()

	userRepository := &repositories.UserRepository{DBConn: DBConn}
	userService := &UserService{IUserRepository: userRepository}
	userController := controllers.UserController{IUserService: userService}

	return userController
}

// GetByID returns a user by id
func (u *UserService) GetByID(id uint) (models.User, error) {
	return u.IUserRepository.GetByID(id)
}

// GetAll returns all users
func (u *UserService) GetAll() ([]models.User, error) {
	users, err := u.IUserRepository.GetAll()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("users not found")
	}

	return users, nil
}
