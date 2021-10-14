package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllAccesses
func GetAllAccesses(c *fiber.Ctx) error {
	db := database.DBConn

	var accesses []models.Access

	if res := db.Joins("User").Joins("Deck").Find(&accesses); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get All accesses",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All accesses",
		Data:    accesses,
		Count:   len(accesses),
	})
}

// GetAllAccesses
func GetAccessesByUserID(c *fiber.Ctx) error {
	db := database.DBConn

	userID := c.Params("userID")

	var accesses []models.Access

	if res := db.Joins("User").Joins("Deck").Where("accesses.user_id = ?", userID).Find(&accesses); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get All accesses",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All accesses",
		Data:    accesses,
		Count:   len(accesses),
	})

}

// GetAccessByID
func GetAccessByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").First(&access, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get access by ID.",
		Data:    *access,
		Count:   1,
	})
}

// GetAccessByUserAndDeckID
func GetAccessByUserAndDeckID(c *fiber.Ctx) error {
	userID := c.Params("userID")
	deckID := c.Params("deckID")

	db := database.DBConn

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", userID, deckID).First(&access).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get access by ID.",
		Data:    *access,
		Count:   1,
	})
}

// POST

// CreateNewAccess
func CreateNewAccess(c *fiber.Ctx) error {
	db := database.DBConn

	access := new(models.Access)

	if err := c.BodyParser(&access); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", access.DeckID).Find(&cards).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	for x := 0; x < len(cards); x++ {
		mem := new(models.Mem)

		if err := db.Where("mems.user_id = ? AND mems.card_id = ?", access.UserID, cards[x].ID).First(&mem).Error; err != nil {
			mem.CardID = cards[x].ID
			mem.UserID = access.UserID
			mem.DeckID = access.DeckID
			mem.Efactor = 2.5
			mem.Interval = 0
			mem.Total = 0
			mem.NextDate = time.Now()
			mem.Quality = 0
			mem.Repetition = 0

			db.Preload("User").Preload("Card").Create(mem)
		}

	}

	db.Preload("User").Preload("Deck").Create(access)

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register an access",
		Data:    *access,
		Count:   1,
	})
}

// PUT

// UpdateAccessByID
func UpdateAccessByID(c *fiber.Ctx) error {
	db := database.DBConn
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

	return c.JSON(models.ResponseHTTP{
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
