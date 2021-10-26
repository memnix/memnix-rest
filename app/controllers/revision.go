package controllers

import (
	"memnixrest/pkg/core"
	"memnixrest/database"
	"memnixrest/app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllRevisions
func GetAllRevisions(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var revisions []models.Revision

	if res := db.Joins("User").Joins("Card").Find(&revisions); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all revisions",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all revisions",
		Data:    revisions,
		Count:   len(revisions),
	})

}

// GetRevisionByID
func GetRevisionByID(c *fiber.Ctx) error {
	db := database.DBConn // Db Conn

	// Params
	id := c.Params("id")

	revision := new(models.Revision)

	if err := db.Joins("User").Joins("Card").First(&revision, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get revision by ID.",
		Data:    *revision,
		Count:   1,
	})
}

// GetRevisionByUserID
func GetRevisionByUserID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("userID")

	revision := new(models.Revision)

	if err := db.Joins("User").Joins("Card").Where("revisions.user_id = ?", id).First(&revision).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get revision by UserID.",
		Data:    *revision,
		Count:   1,
	})
}

// POST

// CreateNewRevision
func CreateNewRevision(c *fiber.Ctx) error {
	db := database.DBConn // Db Conn

	revision := new(models.Revision)

	if err := c.BodyParser(&revision); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Preload("User").Preload("Card").Create(revision)

	mem := core.GetMemByCardAndUser(c, revision.UserID, revision.CardID)
	core.UpdateMem(c, revision, &mem)

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a revision",
		Data:    *revision,
		Count:   1,
	})

}
