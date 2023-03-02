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
}
