package core

import (
	"errors"
	"memnixrest/app/database"
	"memnixrest/app/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GenerateAccess
func GenerateAccess(c *fiber.Ctx, user *models.User, deck *models.Deck) models.ResponseHTTP {
	db := database.DBConn

	if deck.Status != models.DeckPublic && user.Permissions != models.PermAdmin {
		return models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permissions to subscribe to this deck!",
			Data:    nil,
			Count:   0,
		}
	}

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", user.ID, deck.ID).Find(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			access.DeckID = deck.ID
			access.UserID = user.ID
			access.Permission = 1
			db.Preload("User").Preload("Deck").Create(access)
		}

	} else {
		if access.Permission >= models.AccessStudent {
			return models.ResponseHTTP{
				Success: false,
				Message: "You are already subscribed to this deck",
				Data:    nil,
				Count:   0,
			}
		} else {
			access.DeckID = deck.ID
			access.UserID = user.ID
			access.Permission = 1
			db.Preload("User").Preload("Deck").Save(access)
		}
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Success register an access",
		Data:    *access,
		Count:   1,
	}
}
