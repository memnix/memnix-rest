package services

import (
	"bytes"
	"errors"
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"github.com/memnix/memnixrest/interfaces"
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/utils"
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

func (u *UserService) UpdateByID(id uint, newUser *models.User) error {
	user, err := u.IUserRepository.GetByID(id)
	if err != nil {
		return err
	}

	if user.Email != newUser.Email || !bytes.Equal(user.Password, newUser.Password) || user.Permissions != newUser.Permissions {
		return errors.New(utils.ErrorBreak)
	}

	return u.IUserRepository.Update(newUser)

}
