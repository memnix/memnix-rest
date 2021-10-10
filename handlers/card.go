package handlers

import (
	"math/rand"
	"memnixrest/core"
	"memnixrest/database"
	"memnixrest/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetNextCard
func GetNextCard(c *fiber.Ctx) error {
	userIDTemp := c.Params("userID")
	deckIDTemp := c.Params("deckID")

	userID, _ := strconv.Atoi(userIDTemp)
	deckID, _ := strconv.Atoi(deckIDTemp)

	card := core.FetchNextCard(c, uint(userID), uint(deckID))

	//TODO: Handle errors

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get card by ID.",
		Data:    card,
	})

}

// GetRandomDebugCard
func GetRandomDebugCard(c *fiber.Ctx) error {
	rand.Seed(time.Now().UnixNano())
	db := database.DBConn

	var cards []models.Card
	if res := db.Joins("Deck").Find(&cards); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All cards",
			Data:    nil,
		})
	}

	rdm := rand.Intn(len(cards)-0) + 0

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All cards",
		Data:    cards[rdm],
	})
}

// GetAllCards
func GetAllCards(c *fiber.Ctx) error {
	db := database.DBConn

	var cards []models.Card

	if res := db.Joins("Deck").Find(&cards); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All cards",
			Data:    nil,
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All cards",
		Data:    cards,
	})

}

// GetCardByID
func GetCardByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	card := new(models.Card)

	if err := db.Joins("Deck").First(&card, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get card by ID.",
		Data:    *card,
	})
}

// GetCardsFromDeck
func GetCardsFromDeck(c *fiber.Ctx) error {
	id := c.Params("deckID")
	db := database.DBConn

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get cards by ID.",
		Data:    cards,
	})
}

// POST

// CreateNewCard
func CreateNewCard(c *fiber.Ctx) error {
	db := database.DBConn

	card := new(models.Card)

	if err := c.BodyParser(&card); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Preload("Deck").Create(card)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success register a card",
		Data:    *card,
	})
}
