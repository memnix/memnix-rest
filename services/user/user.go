package user

import (
	"context"

	db "github.com/memnix/memnix-rest/db/sqlc"
)

// IUseCase is the user use case interface.
type IUseCase interface {
	// GetName returns the name of the user.
	GetName(ctx context.Context, id string) string
	// GetByID returns the user with the given id.
	GetByID(ctx context.Context, id int32) (db.User, error)
}

// IRepository is the user repository interface.
type IRepository interface {
	// GetName returns the name of the user.
	GetName(ctx context.Context, id int32) string
	// GetByID returns the user with the given id.
	GetByID(ctx context.Context, id int32) (db.User, error)
	// GetByEmail returns the user with the given email.
	GetByEmail(ctx context.Context, email string) (db.User, error)
	// Create creates a new user.
	Create(ctx context.Context, email, password, username string) error
	// Update updates the user with the given id.
	Update(ctx context.Context, id int32, email, password, username string) error
	// Delete deletes the user with the given id.
	Delete(ctx context.Context, id int32) error
	// GetByOauthID returns the user with the given oauth id.
	GetByOauthID(ctx context.Context, id string) (db.User, error)
}

// IRedisRepository is the user redis repository interface.
type IRedisRepository interface {
	// Get returns the value of the given key.
	Get(ctx context.Context, id int32) (string, error)
	// Set sets the value of the given key.
	Set(ctx context.Context, id int32, value string) error
}

type IRistrettoRepository interface {
	// Get returns the value of the given key
	Get(ctx context.Context, id int32) (db.User, error)
	// Set sets the value of the given key
	Set(ctx context.Context, user db.User) error
	// Delete deletes the value of the given key
	Delete(ctx context.Context, id int32) error
}
