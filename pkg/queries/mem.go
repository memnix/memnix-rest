package queries

import (
	"errors"
	"github.com/memnix/memnixrest/pkg/core"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/models"
	"github.com/memnix/memnixrest/pkg/utils"
	"gorm.io/gorm"
	"time"
)

// PostSelfEvaluatedMem updates Mem & MemDate
func PostSelfEvaluatedMem(user *models.User, card *models.Card, quality uint, training bool) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(models.ResponseHTTP)

	memDate := new(models.MemDate)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		user.ID, card.ID).First(&memDate).Error; err != nil {
		res.GenerateError(utils.ErrorRequestFailed) // MemDate not found
		// TODO: Create a default MemDate
		return res
	}

	exMem := FetchMem(memDate.CardID, user.ID)
	if exMem.Efactor == 0 {
		exMem.FillDefaultValues(user.ID, card.ID)
	}

	core.UpdateMemSelfEvaluated(exMem, training, quality)

	res.GenerateSuccess("Success Post Mem", nil, 0)
	return res
}

// PostMem updates MemDate & Mem
func PostMem(user *models.User, card *models.Card, validation *models.CardResponseValidation, training bool) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(models.ResponseHTTP)

	memDate := new(models.MemDate)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		user.ID, card.ID).First(&memDate).Error; err != nil {
		res.GenerateError(utils.ErrorRequestFailed) // MemDate not found
		// TODO: Create a default MemDate
		return res
	}

	exMem := FetchMem(memDate.CardID, user.ID)
	if exMem.Efactor == 0 {
		exMem.FillDefaultValues(user.ID, card.ID)
	}

	if training {
		core.UpdateMemTraining(exMem, validation.Validate)
	} else {
		memDate, err := core.UpdateMem(exMem, validation.Validate)
		if err != nil {
			res.GenerateError(utils.ErrorRequestFailed)
			return res
		}
		t := time.Now()
		if memDate.NextDate.Before(t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour()) * time.Hour)) {
			GetCache().Replace(memDate.UserID, *memDate)
		} else {
			err = GetCache().DeleteItem(memDate.UserID, memDate.ID)
			if err != nil {
				res.GenerateError(utils.ErrorRequestFailed)
				return res
			}
		}
	}

	res.GenerateSuccess("Success Post Mem", nil, 0)
	return res
}

// FetchMem returns last mem of a user on a given card
func FetchMem(cardID, userID uint) *models.Mem {
	db := database.DBConn // DB Conn

	mem := new(models.Mem)
	if err := db.Joins("Card").Where("mems.card_id = ? AND mems.user_id = ?", cardID, userID).Order("id desc").First(&mem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			mem.Efactor = 0
		}
	}
	return mem
}
