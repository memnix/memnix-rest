package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllAccesses
func GetAllAccesses(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var accesses []models.Access // Accesses array

	if res := db.Joins("User").Joins("Deck").Find(&accesses); res.Error != nil {
		// Return error
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all accesses",
			Data:    nil,
			Count:   0,
		})
	}
	// Return success
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all accesses",
		Data:    accesses,
		Count:   len(accesses),
	})
}

// GetAllAccesses
func GetAccessesByUserID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	userID := c.Params("userID")

	var accesses []models.Access // Accesses array

	if res := db.Joins("User").Joins("Deck").Where("accesses.user_id = ?", userID).Find(&accesses); res.Error != nil {
		// Return error
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get accesses of an user",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		// Return success
		Success: true,
		Message: "Get accesses of an user",
		Data:    accesses,
		Count:   len(accesses),
	})
}

// GetAccessByID
func GetAccessByID(c *fiber.Ctx) error {
	db := database.DBConn // Db Conn

	// Params
	id := c.Params("id")

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").First(&access, id).Error; err != nil {
		// Return error
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		// Return success
		Success: true,
		Message: "Success get access by ID.",
		Data:    *access,
		Count:   1,
	})
}

// GetAccessByUserIDAndDeckID
func GetAccessByUserIDAndDeckID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	userID := c.Params("userID")
	deckID := c.Params("deckID")

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", userID, deckID).First(&access).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get access by UserID & DeckID.",
		Data:    *access,
		Count:   1,
	})
}

// GetAllSubAccesses method to get an access
// @Description Get decks a user is sub to
// @Summary get a list of deck
// @Tags Deck
// @Produce json
// @Success 200 {array} models.Access
// @Router /v1/decks/sub [get]
func GetAllSubAccesses(c *fiber.Ctx) error {
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

	var accesses []models.Access // Accesses array

	if err := db.Joins("Deck").Where("accesses.user_id = ? AND accesses.permission >= ?", auth.User.ID, models.AccessStudent).Find(&accesses).Error; err != nil {
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
		Data:    accesses,
		Count:   len(accesses),
	})
}

// POST

// CreateNewAccess
func CreateNewAccess(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	access := new(models.Access)

	if err := c.BodyParser(&access); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Preload("User").Preload("Deck").Create(access)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register an access",
		Data:    *access,
		Count:   1,
	})
}

// PUT

// UpdateAccessByID
func UpdateAccessByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	access := new(models.Access)

	if err := db.First(&access, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if err := UpdateAccess(c, access); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Couldn't update the access",
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update access by Id.",
		Data:    *access,
		Count:   1,
	})
}

// UpdateAccess
func UpdateAccess(c *fiber.Ctx, a *models.Access) error {
	db := database.DBConn

	if err := c.BodyParser(&a); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Preload("User").Preload("Deck").Save(a)

	return nil
}
