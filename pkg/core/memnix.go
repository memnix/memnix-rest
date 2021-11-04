package core

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"time"

	//"time"

	"github.com/gofiber/fiber/v2"
)

// GetUser
func GetUser(userID uint) models.User {
	db := database.DBConn

	user := new(models.User)

	if err := db.First(&user, userID).Error; err != nil {
		return *user
	}

	return *user
}

// FetchNextTodayCard
func FetchNextTodayCard(c *fiber.Ctx, userID uint, deckID uint) models.ResponseHTTP {
	user := GetUser(userID)
	return FetchNextTodayMemByUserAndDeck(c, &user, deckID)
}

// FetchNextCard
func FetchNextCard(c *fiber.Ctx, userID uint, deckID uint) models.ResponseHTTP {
	user := GetUser(userID)
	return FetchNextMemByUserAndDeck(c, &user, deckID)
}

// UpdateMem
func UpdateMem(c *fiber.Ctx, r *models.Mem, validation models.CardResponseValidation) {
	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID = r.UserID
	mem.CardID = r.CardID

	var memType int

	if r.Efactor <= 1.4 || r.Quality <= 1 || r.Repetition < 2 {
		memType = 2
	}

	if validation.Validate {
		if r.Repetition == 0 {
			mem.Interval = 1
		} else if r.Repetition == 1 {
			mem.Interval = 2
		} else if r.Repetition == 2 {
			mem.Interval = 3
		} else {
			mem.Interval = uint((float32(r.Interval) * r.Efactor)) + 1
		}
		mem.Repetition = r.Repetition + 1
		if memType == 2 {
			mem.Quality = 4
		} else {
			mem.Quality = 5
		}
	} else {
		mem.Repetition = 0
		mem.Interval = 0
		if memType == 2 {
			mem.Quality = 1
			if r.Repetition <= 1 {
				mem.Quality = 0
			}
		} else {
			mem.Quality = 3
		}
	}

	mem.Efactor = r.Efactor + (0.1 - (5.0-float32(mem.Quality))*(0.08+(5-float32(mem.Quality)))*0.02)

	if mem.Efactor < 1.3 {
		mem.Efactor = 1.3
	}

	db.Create(mem)

	memDate := new(models.MemDate)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		mem.UserID, mem.CardID).First(&memDate).Error; err != nil {
	}

	memDate.NextDate = time.Now() //time.Now().AddDate(0, 0, int(mem.Interval))

	db.Save(memDate)

	//TODO: Return error
}
