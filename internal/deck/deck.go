package deck

import "github.com/memnix/memnix-rest/domain"

type IUseCase interface {
	// GetByID returns the deck with the given id.
	GetByID(id uint) (domain.Deck, error)
	// Create creates a new deck.
	Create(deck *domain.Deck) error
	// CreateFromUser creates a new deck from the given user.
	CreateFromUser(user domain.User, deck *domain.Deck) error
	// GetByUser returns the decks of the given user.
	GetByUser(user domain.User) ([]domain.Deck, error)
	// GetByLearner returns the decks of the given learner.
	GetByLearner(user domain.User) ([]domain.Deck, error)
	// GetPublic returns the public decks.
	GetPublic() ([]domain.Deck, error)
}

type IRepository interface {
	// GetByID returns the deck with the given id.
	GetByID(id uint) (domain.Deck, error)
	// Create creates a new deck.
	Create(deck *domain.Deck) error
	// Update updates the deck with the given id.
	Update(deck *domain.Deck) error
	// Delete deletes the deck with the given id.
	Delete(id uint) error
	// CreateFromUser creates a new deck from the given user.
	CreateFromUser(user domain.User, deck *domain.Deck) error
	// GetByUser returns the decks of the given user.
	GetByUser(user domain.User) ([]domain.Deck, error)
	// GetByLearner returns the decks of the given learner.
	GetByLearner(user domain.User) ([]domain.Deck, error)
	// GetPublic returns the public decks.
	GetPublic() ([]domain.Deck, error)
}

type IRedisRepository interface {
	// GetByID returns the deck with the given id.
	GetByID(id uint) (string, error)
	// SetByID sets the deck with the given id.
	SetByID(id uint, deck string) error
	// DeleteByID deletes the deck with the given id.
	DeleteByID(id uint) error
	// SetOwnedByUser sets the deck with the given id as owned by the given user.
	SetOwnedByUser(userID uint, decks string) error
	// GetOwnedByUser returns the decks owned by the given user.
	GetOwnedByUser(userID uint) (string, error)
	// DeleteOwnedByUser deletes the decks owned by the given user.
	DeleteOwnedByUser(userID uint) error
	// SetLearningByUser sets the deck with the given id as learnt by the given user.
	SetLearningByUser(userID uint, decks string) error
	// GetLearningByUser returns the decks learnt by the given user.
	GetLearningByUser(userID uint) (string, error)
	// DeleteLearningByUser deletes the decks learnt by the given user.
	DeleteLearningByUser(userID uint) error
}
