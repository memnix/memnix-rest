package repositories

import (
	"github.com/memnix/memnixrest/models"
	"gorm.io/gorm"
)

type CardRepository struct {
	DBConn *gorm.DB
}

func (c *CardRepository) GetAll() ([]models.Card, error) {
	var cards []models.Card
	if res := c.DBConn.Joins("Deck").Find(&cards); res.Error != nil {
		return nil, res.Error
	}

	return cards, nil
}

func (c *CardRepository) GetByID(id uint) (models.Card, error) {
	var card models.Card
	if res := c.DBConn.Joins("Deck").First(&card, id); res.Error != nil {
		return models.Card{}, res.Error
	}

	return card, nil
}

func (c *CardRepository) GetCardsFromDeck(deckID uint) ([]models.Card, error) {
	var cards []models.Card
	if res := c.DBConn.Joins("Deck").Where("deck_id = ?", deckID).Find(&cards); res.Error != nil {
		return nil, res.Error
	}

	return cards, nil
}

func (c *CardRepository) CreateCard(card models.Card) error {
	if res := c.DBConn.Create(&card); res.Error != nil {
		return res.Error
	}

	return nil
}
