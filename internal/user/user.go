package user

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
)

// IUseCase is the user use case interface.
type IUseCase interface {
	// GetName returns the name of the user.
	GetName(ctx context.Context, id string) string
	// GetByID returns the user with the given id.
	GetByID(ctx context.Context, id uint) (domain.User, error)
}

// IRepository is the user repository interface.
type IRepository interface {
	// GetName returns the name of the user.
	GetName(ctx context.Context, id uint) string
	// GetByID returns the user with the given id.
	GetByID(ctx context.Context, id uint) (domain.User, error)
	// GetByEmail returns the user with the given email.
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	// Create creates a new user.
	Create(ctx context.Context, user *domain.User) error
	// Update updates the user with the given id.
	Update(ctx context.Context, user *domain.User) error
	// Delete deletes the user with the given id.
	Delete(ctx context.Context, id uint) error
	// GetByOauthID returns the user with the given oauth id.
	GetByOauthID(ctx context.Context, id string) (domain.User, error)
}

// IRedisRepository is the user redis repository interface.
type IRedisRepository interface {
	// Get returns the value of the given key.
	Get(ctx context.Context, id uint) (string, error)
	// Set sets the value of the given key.
	Set(ctx context.Context, id uint, value string) error
}
