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

	res := core.FetchNextCard(c, uint(userID), uint(deckID))

	if !res.Success {
		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Next card not found",
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get card by ID.",
		Data:    res.Data,
		Count:   1,
	})

}

// GetTodayNextCard
func GetTodayNextCard(c *fiber.Ctx) error {
	userIDTemp := c.Params("userID")
	deckIDTemp := c.Params("deckID")

	userID, _ := strconv.Atoi(userIDTemp)
	deckID, _ := strconv.Atoi(deckIDTemp)

	res := core.FetchNextTodayCard(c, uint(userID), uint(deckID))
	if !res.Success {
		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "No more card for today!",
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get card Today's card.",
		Data:    res.Data,
		Count:   1,
	})
}

// GetRandomDebugCard
func GetRandomDebugCard(c *fiber.Ctx) error {
	rand.Seed(time.Now().UnixNano())
	db := database.DBConn

	var cards []models.Card
	if res := db.Joins("Deck").Find(&cards); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get All cards",
			Data:    nil,
			Count:   0,
		})
	}

	rdm := rand.Intn(len(cards)-0) + 0

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All cards",
		Data:    cards[rdm],
		Count:   len(cards),
	})
}

// GetAllCards
func GetAllCards(c *fiber.Ctx) error {
	db := database.DBConn

	var cards []models.Card

	if res := db.Joins("Deck").Find(&cards); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get All cards",
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
	id := c.Params("id")
	db := database.DBConn

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
	id := c.Params("deckID")
	db := database.DBConn

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
		Message: "Success get cards by ID.",
		Data:    cards,
		Count:   len(cards),
	})
}

// POST

// CreateNewCard
func CreateNewCard(c *fiber.Ctx) error {
	db := database.DBConn

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
