package user

import "github.com/memnix/memnix-rest/domain"

type IUseCase interface {
	// GetName returns the name of the user.
	GetName(id string) string
	// GetByID returns the user with the given id.
	GetByID(id uint) (domain.User, error)
}

type IRepository interface {
	// GetName returns the name of the user.
	GetName(id uint) string
	// GetByID returns the user with the given id.
	GetByID(id uint) (domain.User, error)
	// GetByEmail returns the user with the given email.
	GetByEmail(email string) (domain.User, error)
	// Create creates a new user.
	Create(user *domain.User) error
	// Update updates the user with the given id.
	Update(user *domain.User) error
	// Delete deletes the user with the given id.
	Delete(id uint) error
	// GetByOauthID returns the user with the given oauth id.
	GetByOauthID(id string) (domain.User, error)
}
