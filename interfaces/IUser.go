package interfaces

import "github.com/memnix/memnixrest/models"

type IUserRepository interface {
	// GetByID returns a user by id
	GetByID(id uint) (models.User, error)
	// GetAll returns all users
	GetAll() ([]models.User, error)
}

type IUserService interface {
	// GetByID returns a user by id
	GetByID(id uint) (models.User, error)
	// GetAll returns all users
	GetAll() ([]models.User, error)
}
