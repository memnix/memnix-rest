package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllIdentifiers
func GetAllIdentifiers(c *fiber.Ctx) error {
	db := database.DBConn

	var identifiers []models.Identifier

	if res := db.Find(&identifiers); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All users",
			Data:    nil,
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All users",
		Data:    identifiers,
	})

}

// GetIdentifierByID
func GetIdentifierByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	identifier := new(models.Identifier)

	if err := db.First(&identifier, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get identifier by ID.",
		Data:    *identifier,
	})
}

// GetIdentifierByDiscordID
func GetIdentifierByDiscordID(c *fiber.Ctx) error {
	id := c.Params("discordID")
	db := database.DBConn

	identifier := new(models.Identifier)

	if err := db.Where("identifier.discord_id = ?", id).First(&identifier).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get identifier by ID.",
		Data:    *identifier,
	})
}

// GetIdentifierByUserID
func GetIdentifierByUserID(c *fiber.Ctx) error {
	id := c.Params("userID")
	db := database.DBConn

	identifier := new(models.Identifier)

	if err := db.Where("identifier.user_id = ?", id).First(&identifier).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get identifier by ID.",
		Data:    *identifier,
	})
}

// POST

// CreateNewUser
func CreateNewIdentifier(c *fiber.Ctx) error {
	db := database.DBConn

	identifier := new(models.Identifier)

	if err := c.BodyParser(&identifier); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Create(identifier)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success register an user",
		Data:    *identifier,
	})
}
