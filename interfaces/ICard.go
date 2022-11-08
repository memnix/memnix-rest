package interfaces

import "github.com/memnix/memnixrest/models"

type ICardRepository interface {
	// GetAll returns all cards
	GetAll() ([]models.Card, error)
	// GetByID returns a card by id
	GetByID(id uint) (models.Card, error)
	// GetCardsFromDeck returns all cards from a deck
	GetCardsFromDeck(deckID uint) ([]models.Card, error)
	// CreateCard creates a card
	CreateCard(card models.Card) error
}

type ICardService interface {
	// GetAll returns all cards
	GetAll() ([]models.Card, error)
	// GetByID returns a card by id
	GetByID(id uint) (models.Card, error)
	// GetCardsFromDeck returns all cards from a deck
	GetCardsFromDeck(deckID uint) ([]models.Card, error)
	// CreateCard creates a card
	CreateCard(card models.Card) error
}
