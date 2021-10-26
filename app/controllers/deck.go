package controllers

import (
	"memnixrest/app/models"
	"memnixrest/app/database"
	"memnixrest/pkg/core"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllDecks
func GetAllDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var decks []models.Deck

	if res := db.Find(&decks); res.Error != nil {

		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all decks",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all decks",
		Data:    decks,
		Count:   len(decks),
	})

}

// GetDeckByID
func GetDeckByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	deck := new(models.Deck)

	if err := db.First(&deck, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get deck by ID.",
		Data:    *deck,
		Count:   1,
	})
}

// GetAllSubDecks
func GetAllSubDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("userID")

	var decks []models.Deck

	if err := db.Joins("JOIN accesses ON accesses.deck_id = decks.id AND accesses.user_id = ? AND accesses.permission > 0", id).Find(&decks).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all sub decks",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all sub decks",
		Data:    decks,
		Count:   len(decks),
	})
}

// GetAllPublicDecks
func GetAllPublicDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var decks []models.Deck

	if err := db.Where("decks.status = 2").Find(&decks).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all public decks",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all public decks",
		Data:    decks,
		Count:   len(decks),
	})
}

// POST

// CreateNewDeck
func CreateNewDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	deck := new(models.Deck)

	if err := c.BodyParser(&deck); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Create(deck)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a deck",
		Data:    *deck,
		Count:   1,
	})
}

// UnsubToDeck
func UnSubToDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	deckID := c.Params("deckID")
	userID := c.Params("userID")

	access := new(models.Access)
	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", userID, deckID).Find(&access).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "This user isn't sub to the deck",
			Data:    nil,
			Count:   0,
		})
	}

	access.Permission = 0
	db.Preload("User").Preload("Deck").Save(access)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success unsub to the deck",
		Data:    nil,
		Count:   0,
	})
}

// SubToDeck
func SubToDeck(c *fiber.Ctx) error {
	deckID := c.Params("deckID")
	userID := c.Params("userID")

	userID_temp, _ := strconv.Atoi(userID)
	deckID_temp, _ := strconv.Atoi(deckID)

	db := database.DBConn

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", deckID).Find(&cards).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if err := core.GenerateAccess(c, uint(userID_temp), uint(deckID_temp)); !err.Success {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Couldn't generate access !",
			Data:    nil,
			Count:   0,
		})
	}

	if err := core.GenerateMem(c, uint(userID_temp), uint(deckID_temp)); !err.Success {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Couldn't generate mems !",
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success subscribing to deck",
		Data:    nil,
		Count:   0,
	})
}
