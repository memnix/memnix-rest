package core

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ComputeInterval function
func ComputeInterval(r *models.Mem) uint {
	switch r.Repetition {
	case 0:
		return 1
	case 1, 2:
		return 2
	case 3:
		return 3
	default:
		return uint(float32(r.Interval)*r.Efactor*0.75) + 1
	}
}

func ComputeQualitySuccess(memType int, repetition uint) uint {

	if memType == 2 {
		return 3
	} else {
		if repetition > 3 {
			return 5
		}
		return 4
	}
}

func ComputeQualityFailed(memType int, repetition uint) uint {
	if memType == 2 {
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

func ComputeEfactor(r *models.Mem) float32 {
	eFactor := r.Efactor + (0.1 - (5.0-float32(r.Quality))*(0.08+(5-float32(r.Quality)))*0.02)

	if eFactor < 1.3 {
		return 1.3
	}
	return eFactor
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

// UpdateMem function
func UpdateMem(_ *fiber.Ctx, r *models.Mem, validation models.CardResponseValidation) {
	//TODO: Rewrite functions

	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	var memType int

	if r.Efactor <= 2 || r.Repetition < 2 || (r.Efactor <= 2.3 && r.Repetition < 4) {
		memType = 2
	}

	if validation.Validate {
		mem.Interval = ComputeInterval(r)
		mem.Repetition = r.Repetition + 1
		r.Quality = ComputeQualitySuccess(memType, r.Repetition)

	} else {
		mem.Repetition = 0
		mem.Interval = 0
		r.Quality = ComputeQualityFailed(memType, r.Repetition)
	}

	mem.Efactor = ComputeEfactor(r)

	db.Save(r)
	db.Create(mem)

	UpdateMemDate(mem)

}
