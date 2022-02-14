package core

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// UpdateMem function
func UpdateMem(_ *fiber.Ctx, r *models.Mem, validation models.CardResponseValidation) {
	//TODO: Rewrite functions

	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID = r.UserID
	mem.CardID = r.CardID

	var memType int

	if r.Efactor <= 2 || r.Repetition < 2 || (r.Efactor <= 2.3 && r.Repetition < 4) {
		memType = 2
	}

	if validation.Validate {
		if r.Repetition == 0 {
			mem.Interval = 1
		} else if r.Repetition == 1 {
			mem.Interval = 2
		} else if r.Repetition == 2 {
			mem.Interval = 2
		} else if r.Repetition == 3 {
			mem.Interval = 3
		} else {
			mem.Interval = uint(float32(r.Interval)*r.Efactor*0.75) + 1
		}
		mem.Repetition = r.Repetition + 1
		if memType == 2 {
			r.Quality = 3
		} else {
			r.Quality = 4
			if r.Repetition > 3 {
				r.Quality = 5
			}
		}
	} else {
		mem.Repetition = 0
		mem.Interval = 0
		if memType == 2 {
			r.Quality = 1
			if r.Repetition <= 1 {
				r.Quality = 0
			}
		} else {
			r.Quality = 2
		}
	}
	mem.Efactor = r.Efactor + (0.1 - (5.0-float32(r.Quality))*(0.08+(5-float32(r.Quality)))*0.02)

	if mem.Efactor < 1.3 {
		mem.Efactor = 1.3
	}

	db.Save(r)
	db.Create(mem)

	memDate := new(models.MemDate)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		mem.UserID, mem.CardID).First(&memDate).Error; err != nil {
	}

	memDate.NextDate = time.Now().AddDate(0, 0, int(mem.Interval))

	db.Save(memDate)

	//TODO: Return error
}
