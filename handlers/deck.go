package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllDecks
func GetAllDecks(c *fiber.Ctx) error {
	db := database.DBConn

	var decks []models.Deck

	if res := db.Find(&decks); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get All decks",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All decks",
		Data:    decks,
		Count:   len(decks),
	})

}

// GetDeckByID
func GetDeckByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	deck := new(models.Deck)

	if err := db.First(&deck, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get deck by ID.",
		Data:    *deck,
		Count:   1,
	})
}

// POST

// CreateNewDeck
func CreateNewDeck(c *fiber.Ctx) error {
	db := database.DBConn

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

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a deck",
		Data:    *deck,
		Count:   1,
	})
}
