package services

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"github.com/memnix/memnixrest/interfaces"
	"github.com/memnix/memnixrest/models"
)

type CardService struct {
	interfaces.ICardRepository
}

func (k *kernel) InjectCardController() controllers.CardController {
	DBConn := infrastructures.GetDBConn()

	cardRepository := &repositories.CardRepository{DBConn: DBConn}
	cardService := &CardService{ICardRepository: cardRepository}
	cardController := controllers.CardController{ICardService: cardService}

	return cardController
}

func (c *CardService) GetAll() ([]models.Card, error) {
	return c.ICardRepository.GetAll()
}

func (c *CardService) GetByID(id uint) (models.Card, error) {
	return c.ICardRepository.GetByID(id)
}

func (c *CardService) GetCardsFromDeck(deckID uint) ([]models.Card, error) {
	cards, err := c.ICardRepository.GetCardsFromDeck(deckID)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *CardService) CreateCard(card models.Card) error {
	return c.ICardRepository.CreateCard(card)
}
