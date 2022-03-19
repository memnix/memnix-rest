package core

import (
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/pkg/database"
)

// UpdateMemDate computes NextDate and set it
func UpdateMemDate(mem *models.Mem) {
	db := database.DBConn
	memDate := new(models.MemDate)

	_ = db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		mem.UserID, mem.CardID).First(&memDate).Error
	//TODO: Error handling

	memDate.ComputeNextDate(int(mem.Interval))

	db.Save(memDate)

	//TODO: Return error
}

// UpdateMemTraining computes and set mem values
func UpdateMemTraining(r *models.Mem, validation bool) {
	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	if validation {
		r.ComputeQualitySuccess()
	} else {
		r.ComputeQualityFail()
	}

	mem.ComputeTrainingEfactor(r.Efactor, r.Quality)
	mem.Interval, mem.Repetition = r.Interval, r.Repetition

	db.Save(r)
	db.Create(mem)
}

// UpdateMem computes and set mem values
func UpdateMem(r *models.Mem, validation bool) {

	db := database.DBConn

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	if validation {
		mem.ComputeInterval(r.Interval, r.Efactor, r.Repetition)
		mem.Repetition = r.Repetition + 1
		r.ComputeQualitySuccess()

	} else {
		mem.Repetition = 0
		mem.Interval = 0
		r.ComputeQualityFail()
	}

	mem.Quality = models.MemQualityNone

	mem.ComputeEfactor(r.Efactor, r.Quality)

	db.Save(r)
	db.Create(mem)

	UpdateMemDate(mem)

}
