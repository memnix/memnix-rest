package core

import (
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/models"
	"strings"
)

// UpdateMemSelfEvaluated computes self evaluated mem
func UpdateMemSelfEvaluated(r *models.Mem, training bool, quality uint) {
	db := infrastructures.GetDBConn()

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	mem.Quality = models.MemQualityNone
	r.Quality = models.MemQuality(quality)

	if training {
		mem.ComputeTrainingEfactor(r.Efactor, r.Quality)
	} else {
		mem.ComputeEfactor(r.Efactor, r.Quality)
	}

	mem.Interval, mem.Repetition = r.Interval, r.Repetition

	db.Save(r)
	db.Create(mem)
}

// UpdateMemDate computes NextDate and set it
func UpdateMemDate(mem *models.Mem) (*models.MemDate, error) {
	db := infrastructures.GetDBConn()
	memDate := new(models.MemDate)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		mem.UserID, mem.CardID).First(&memDate).Error; err != nil {
		return nil, err
	}

	memDate.ComputeNextDate(int(mem.Interval))

	memDate.LearningStage = mem.LearningStage

	db.Save(memDate)

	return memDate, nil
}

// UpdateMemTraining computes and set mem values
func UpdateMemTraining(r *models.Mem, validation bool) {
	db := infrastructures.GetDBConn()

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	if validation {
		r.ComputeQualitySuccess()
	} else {
		r.ComputeQualityFail()
	}

	mem.Quality = models.MemQualityNone

	mem.ComputeTrainingEfactor(r.Efactor, r.Quality)
	mem.Interval, mem.Repetition = r.Interval, r.Repetition

	db.Save(r)
	db.Create(mem)
}

// UpdateMem computes and set mem values
func UpdateMem(r *models.Mem, validation bool) (*models.MemDate, error) {
	db := infrastructures.GetDBConn()

	mem := new(models.Mem)

	mem.UserID, mem.CardID = r.UserID, r.CardID

	if validation {
		mem.ComputeInterval(r.Interval, r.Efactor, r.Repetition)
		mem.Repetition = r.Repetition + 1
		mem.ComputeLearningStage()
		r.ComputeQualitySuccess()
	} else {
		mem.Repetition = 0
		mem.Interval = 0
		mem.LearningStage = models.StageToLearn
		r.ComputeQualityFail()
	}

	mem.Quality = models.MemQualityNone

	mem.ComputeEfactor(r.Efactor, r.Quality)

	db.Save(r)
	db.Create(mem)

	memDate, err := UpdateMemDate(mem)
	if err != nil {
		return nil, err
	}

	return memDate, nil
}

func ValidateAnswer(response string, card *models.Card) bool {
	var respString, answerString string
	if card.Spaces {
		respString = strings.Join(strings.Fields(response), " ")
		answerString = strings.Join(strings.Fields(card.Answer), " ")
	} else {
		respString = strings.ReplaceAll(response, " ", "")
		answerString = strings.ReplaceAll(card.Answer, " ", "")
	}
	if card.Case {
		return strings.Compare(respString, answerString) == 0
	}
	return strings.EqualFold(respString, answerString)
}
