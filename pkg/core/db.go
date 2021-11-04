package core

import (
	"errors"
	"memnixrest/app/database"
	"memnixrest/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// FetchNextTodayMemByUserAndDeck
func FetchNextTodayMemByUserAndDeck(c *fiber.Ctx, user *models.User, deck_id uint) models.ResponseHTTP {
	db := database.DBConn
	mem := new(models.Mem)
	// Get next card with date condition
	t := time.Now()
	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mems.user_id = ? AND mems.deck_id =? AND mems.next_date < ?",
		&user.ID, deck_id, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour)).Limit(1).Order("next_date asc").Find(&mem).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next today mem not found",
			Data:    nil,
		}
	}
	return models.ResponseHTTP{
		Success: true,
		Message: "Get Next Today Mem",
		Data:    *mem,
	}
}

// FetchNextMemByUserAndDeck
func FetchNextMemByUserAndDeck(c *fiber.Ctx, user *models.User, deck_id uint) models.ResponseHTTP {
	db := database.DBConn
	mem := new(models.Mem)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mems.user_id = ? AND mems.deck_id =?", &user.ID, deck_id).Limit(1).Order("next_date asc").Find(&mem).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next mem not found",
			Data:    nil,
		}

	}
	return models.ResponseHTTP{
		Success: true,
		Message: "Get Next Mem",
		Data:    *mem,
	}
}

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
