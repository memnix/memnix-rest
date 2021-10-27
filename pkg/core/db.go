package core

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
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

// GenerateAccess
func GenerateAccess(c *fiber.Ctx, userID uint, deckID uint) models.ResponseHTTP {
	db := database.DBConn

	access := new(models.Access)
	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", userID, deckID).Find(&access).Error; err != nil {
		access.DeckID = deckID
		access.UserID = userID
		access.Permission = 1
		db.Preload("User").Preload("Deck").Create(access)
	} else {
		access.DeckID = deckID
		access.UserID = userID
		access.Permission = 1
		db.Preload("User").Preload("Deck").Save(access)
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Success register an access",
		Data:    *access,
		Count:   1,
	}
}

// GenerateMem
func GenerateMem(c *fiber.Ctx, userID uint, deckID uint) models.ResponseHTTP {
	db := database.DBConn
	var cards []models.Card // Cards array

	if err := db.Joins("Deck").Where("cards.deck_id = ?", deckID).Find(&cards).Error; err != nil {
		return (models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	for x := 0; x < len(cards); x++ {
		mem := new(models.Mem)

		if err := db.Where("mems.user_id = ? AND mems.card_id = ?", userID, cards[x].ID).First(&mem).Error; err != nil {
			mem.CardID = cards[x].ID
			mem.UserID = userID
			mem.Efactor = 2.5
			mem.Interval = 0
			mem.Total = 0
			mem.NextDate = time.Now()
			mem.Quality = 0
			mem.Repetition = 0

			db.Preload("User").Preload("Card").Create(mem)
		}
	}

	return (models.ResponseHTTP{
		Success: true,
		Message: "Generate mems",
		Data:    nil,
		Count:   0,
	})
}
