package core

import (
	"memnixrest/app/models"
	"memnixrest/pkg/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ComputeQualitySuccess(memType models.CardType, repetition uint) uint {

	if memType == models.CardMCQ {
		return 3
	} else {
		if repetition > 3 {
			return 5
		}
		return 4
	}
}

func ComputeQualityFailed(memType models.CardType, repetition uint) uint {
	if memType == models.CardMCQ {
		if repetition <= 3 {
			return 0
		}
		return 1
	}
	if repetition <= 4 {
		return 1
	}
	return 2

}

func UpdateMemDate(mem *models.Mem) {
	db := database.DBConn

	memDate := new(models.MemDate)

	_ = db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		mem.UserID, mem.CardID).First(&memDate).Error

	memDate.NextDate = time.Now().AddDate(0, 0, int(mem.Interval))

	db.Save(memDate)

	//TODO: Return error

}

func UpdateMemTraining(r *models.Mem, validation *models.CardResponseValidation) {
	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	memType := r.GetMemType()

	if validation.Validate {
		r.Quality = ComputeQualitySuccess(memType, r.Repetition)
	} else {
		r.Quality = ComputeQualityFailed(memType, r.Repetition)
	}

	mem.ComputeTrainingEfactor(r.Efactor, r.Quality)
	mem.Interval = r.Interval
	mem.Repetition = r.Repetition

	db.Save(r)
	db.Create(mem)
}

// UpdateMem function
func UpdateMem(_ *fiber.Ctx, r *models.Mem, validation *models.CardResponseValidation) {

	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	memType := r.GetMemType()

	if validation.Validate {
		mem.ComputeInterval(r.Interval, r.Efactor, r.Repetition)
		mem.Repetition = r.Repetition + 1
		r.Quality = ComputeQualitySuccess(memType, r.Repetition)

	} else {
		mem.Repetition = 0
		mem.Interval = 0

		r.Quality = ComputeQualityFailed(memType, r.Repetition)
	}

	mem.ComputeEfactor(r.Efactor, r.Quality)

	db.Save(r)
	db.Create(mem)

	UpdateMemDate(mem)

}
