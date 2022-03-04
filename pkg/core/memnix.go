package core

import (
	"memnixrest/app/models"
	"memnixrest/pkg/database"
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

// UpdateMem function
func UpdateMem(_ *fiber.Ctx, r *models.Mem, validation *models.CardResponseValidation, training bool) {

	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	memType := r.GetMemType()

	if validation.Validate {
		if !training {
			mem.Interval = ComputeInterval(r)
		}
		mem.Repetition = r.Repetition + 1
		r.Quality = ComputeQualitySuccess(memType, r.Repetition)

	} else {
		mem.Repetition = 0
		if !training {
			mem.Interval = 0
		}
		r.Quality = ComputeQualityFailed(memType, r.Repetition)
	}

	if training {
		mem.ComputeTrainingEfactor(r.Efactor, r.Quality)
	} else {
		mem.ComputeEfactor(r.Efactor, r.Quality)
	}

	db.Save(r)
	db.Create(mem)

	if !training {
		UpdateMemDate(mem)
	}
}
