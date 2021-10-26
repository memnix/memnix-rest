package handlers

import (
	"math/rand"
	"memnixrest/app/models"
	"memnixrest/database"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET

// DEBUG: GetRandomDebugCard
func GetRandomDebugCard(c *fiber.Ctx) error {
	rand.Seed(time.Now().UnixNano()) // Random seed

	db := database.DBConn // DB Conn

	var cards []models.Card
	if res := db.Joins("Deck").Find(&cards); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get a random card",
			Data:    nil,
			Count:   0,
		})
	}

	rdm := rand.Intn(len(cards)-0) + 0

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get a random card",
		Data:    cards[rdm],
		Count:   len(cards),
	})
}

// GetAllCards
func GetAllCards(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var cards []models.Card

	if res := db.Joins("Deck").Find(&cards); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all cards",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All cards",
		Data:    cards,
		Count:   len(cards),
	})

}

// GetCardByID
func GetCardByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	card := new(models.Card)

	if err := db.Joins("Deck").First(&card, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get card by ID.",
		Data:    *card,
		Count:   1,
	})
}

// GetCardsFromDeck
func GetCardsFromDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("deckID")

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get cards from deck.",
		Data:    cards,
		Count:   len(cards),
	})
}

// POST

// CreateNewCard
func CreateNewCard(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	card := new(models.Card)

	if err := c.BodyParser(&card); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Preload("Deck").Create(card)

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a card",
		Data:    *card,
		Count:   1,
	})
}
