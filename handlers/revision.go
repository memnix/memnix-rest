package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllRevisions
func GetAllRevisions(c *fiber.Ctx) error {
	db := database.DBConn

	var revisions []models.Revision

	if res := db.Joins("User").Joins("Card").Find(&revisions); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All revisions",
			Data:    nil,
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All revisions",
		Data:    revisions,
	})

}

// GetRevisionByID
func GetRevisionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	revision := new(models.Revision)

	if err := db.Joins("User").Joins("Card").First(&revision, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get revision by ID.",
		Data:    *revision,
	})
}

// GetRevisionByUserID
func GetRevisionByUserID(c *fiber.Ctx) error {
	id := c.Params("userID")
	db := database.DBConn

	revision := new(models.Revision)

	if err := db.Joins("User").Joins("Card").Where("revisions.user_id = ?", id).First(&revision).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get revision by ID.",
		Data:    *revision,
	})
}

// GetRevisionByCardID
func GetRevisionByCardID(c *fiber.Ctx) error {
	id := c.Params("userID")
	db := database.DBConn

	revision := new(models.Revision)

	if err := db.Joins("User").Joins("Card").Where("revisions.card_id = ?", id).First(&revision).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get revision by ID.",
		Data:    *revision,
	})
}

// POST

// CreateNewRevision
func CreateNewRevision(c *fiber.Ctx) error {
	db := database.DBConn

	revision := new(models.Revision)

	if err := c.BodyParser(&revision); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Preload("User").Preload("Card").Create(revision)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success register a revision",
		Data:    *revision,
	})
}
