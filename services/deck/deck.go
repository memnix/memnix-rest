package deck

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
)

// IUseCase is the interface for the deck use case.
type IUseCase interface {
	// GetByID returns the deck with the given id.
	GetByID(ctx context.Context, id uint) (domain.Deck, error)
	// Create creates a new deck.
	Create(ctx context.Context, deck *domain.Deck) error
	// CreateFromUser creates a new deck from the given user.
	CreateFromUser(ctx context.Context, user domain.User, deck *domain.Deck) error
	// GetByUser returns the decks of the given user.
	GetByUser(ctx context.Context, user domain.User) ([]domain.Deck, error)
	// GetByLearner returns the decks of the given learner.
	GetByLearner(ctx context.Context, user domain.User) ([]domain.Deck, error)
	// GetPublic returns the public decks.
	GetPublic(ctx context.Context) ([]domain.Deck, error)
}

// IRepository is the interface for the deck repository.
type IRepository interface {
	// GetByID returns the deck with the given id.
	GetByID(ctx context.Context, id uint) (domain.Deck, error)
	// Create creates a new deck.
	Create(ctx context.Context, deck *domain.Deck) error
	// Update updates the deck with the given id.
	Update(ctx context.Context, deck *domain.Deck) error
	// Delete deletes the deck with the given id.
	Delete(ctx context.Context, id uint) error
	// CreateFromUser creates a new deck from the given user.
	CreateFromUser(ctx context.Context, user domain.User, deck *domain.Deck) error
	// GetByUser returns the decks of the given user.
	GetByUser(ctx context.Context, user domain.User) ([]domain.Deck, error)
	// GetByLearner returns the decks of the given learner.
	GetByLearner(ctx context.Context, user domain.User) ([]domain.Deck, error)
	// GetPublic returns the public decks.
	GetPublic(ctx context.Context) ([]domain.Deck, error)
}

// IRedisRepository is the interface for the deck redis repository.
type IRedisRepository interface {
	// GetByID returns the deck with the given id.
	GetByID(ctx context.Context, id uint) (string, error)
	// SetByID sets the deck with the given id.
	SetByID(ctx context.Context, id uint, deck string) error
	// DeleteByID deletes the deck with the given id.
	DeleteByID(ctx context.Context, id uint) error
	// SetOwnedByUser sets the deck with the given id as owned by the given user.
	SetOwnedByUser(ctx context.Context, userID uint, decks string) error
	// GetOwnedByUser returns the decks owned by the given user.
	GetOwnedByUser(ctx context.Context, userID uint) (string, error)
	// DeleteOwnedByUser deletes the decks owned by the given user.
	DeleteOwnedByUser(ctx context.Context, userID uint) error
	// SetLearningByUser sets the deck with the given id as learnt by the given user.
	SetLearningByUser(ctx context.Context, userID uint, decks string) error
	// GetLearningByUser returns the decks learnt by the given user.
	GetLearningByUser(ctx context.Context, userID uint) (string, error)
	// DeleteLearningByUser deletes the decks learnt by the given user.
	DeleteLearningByUser(ctx context.Context, userID uint) error
}
