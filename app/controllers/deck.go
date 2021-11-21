package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"memnixrest/pkg/queries"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllDecks method to get all decks
// @Description Get every deck. Shouldn't really be used, consider using /v1/decks/public instead !
// @Summary get all decks
// @Tags Deck
// @Produce json
// @Success 200 {object} models.Deck
// @Router /v1/decks [get]
func GetAllDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

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

// GetDeckByID method to get a deck
// @Description Get a deck by tech ID
// @Summary get a deck
// @Tags Deck
// @Produce json
// @Param id path int true "Deck ID"
// @Success 200 {model} models.Deck
// @Router /v1/decks/{id} [get]
func GetDeckByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

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

// GetAllSubDecks method to get a deck
// @Description Get decks a user is sub to
// @Summary get a list of deck
// @Tags Deck
// @Produce json
// @Param userID path int true "user ID"
// @Success 200 {array} models.Deck
// @Router /v1/decks/user/{userID} [get]
func GetAllSubDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("userID")

	var decks []models.Deck

	if err := db.Joins("JOIN accesses ON accesses.deck_id = decks.id AND accesses.user_id = ? AND accesses.permission >= ?", id, models.AccessStudent).Find(&decks).Error; err != nil {
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

// GetAllPublicDecks method to get a deck
// @Description Get all public deck
// @Summary get a list of deck
// @Tags Deck
// @Produce json
// @Success 200 {model} models.Deck
// @Router /v1/decks/public [get]
func GetAllPublicDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var decks []models.Deck

	if err := db.Where("decks.status = ?", models.DeckPublic).Find(&decks).Error; err != nil {
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

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

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

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

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
	db := database.DBConn // DB Conn
	deckID := c.Params("deckID")
	deck := new(models.Deck)

	// Check auth
	auth := CheckAuth(c, models.PermUser)
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	if err := db.First(&deck, deckID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if err := queries.GenerateAccess(c, &auth.User, deck); !err.Success {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Message,
			Data:    nil,
			Count:   0,
		})
	}

	if err := queries.PopulateMemDate(c, &auth.User, deck); !err.Success {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Message,
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
