package core

import (
	"memnixrest/database"
	"memnixrest/models"

	"github.com/gofiber/fiber/v2"
)

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

func FetchNextMemByUserAndDeck(c *fiber.Ctx, user *models.User, deck_id uint) models.Mem {
	db := database.DBConn
	mem := new(models.Mem)
	// Get next card with date condition
	// 	t := time.Now()
	// "mem.next_date < ?",  t.AddDate(0, 0, 1).Add(time.Duration(-t.Hour()) * time.Hour

	if err := db.Joins("Card").Where("mems.user_id = ? AND mems.deck_id =?", &user.ID, deck_id).Limit(1).Order("next_date asc").Find(&mem).Error; err != nil {
		return *mem

		// TODO: handle error
	}

	return *mem

}
