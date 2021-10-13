package core

import (
	"memnixrest/database"
	"memnixrest/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetMemByID
func GetMemByID(c *fiber.Ctx, id uint) models.Mem {
	db := database.DBConn
	mem := new(models.Mem)

	if err := db.Joins("User").Joins("Card").First(&mem, id).Error; err != nil {
		return *mem
	}
	return *mem
}

// GetMemByCardAndUser
func GetMemByCardAndUser(c *fiber.Ctx, userID uint, cardID uint) models.Mem {

	db := database.DBConn

	mem := new(models.Mem)

	if err := db.Joins("User").Joins("Card").Where("mems.user_id = ? AND mems.card_id = ?", userID, cardID).First(&mem).Error; err != nil {
		return *mem
		// TODO: Handle errors
	}
	return *mem
}

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
