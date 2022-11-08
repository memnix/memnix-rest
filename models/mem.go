package models

import (
	"gorm.io/gorm"
)

// Mem structure
type Mem struct {
	gorm.Model    `swaggerignore:"true"`
	UserID        uint          `json:"user_id" example:"1"`
	User          User          `swaggerignore:"true"`
	CardID        uint          `json:"card_id" example:"1"`
	Card          Card          `swaggerignore:"true"`
	Quality       MemQuality    `json:"quality" example:"0"`
	Repetition    uint          `json:"repetition" example:"0" `
	Efactor       float32       `json:"e_factor" example:"2.5"`
	Interval      uint          `json:"interval" example:"0"`
	LearningStage LearningStage `json:"learning_stage"`
}

// MemQuality enum type
type MemQuality int64

const (
	MemQualityNone     MemQuality = -1
	MemQualityBlackout MemQuality = iota
	MemQualityErrorMCQ
	MemQualityErrorHints
	MemQualityError
	MemQualityGoodMCQ
	MemQualityPerfect
)

// FillDefaultValues to fill a Mem with default values for given UserID and CardID
func (mem *Mem) FillDefaultValues(userID, cardID uint) {
	mem.UserID = userID
	mem.CardID = cardID
	mem.Quality = MemQualityBlackout
	mem.Repetition = 0
	mem.Efactor = 2.5
	mem.Interval = 0
	mem.LearningStage = StageToLearn
}

// ComputeEfactor calculates and sets new efactor using oldEfactor and MemQuality
func (mem *Mem) ComputeEfactor(oldEfactor float32, quality MemQuality) {
	eFactor := oldEfactor + (0.1 - (5.0-float32(quality))*(0.08+(5-float32(quality)))*0.02)

	if eFactor < 1.3 {
		mem.Efactor = 1.3
	} else {
		mem.Efactor = eFactor
	}
}

// ComputeTrainingEfactor calculates and sets new efactor using oldEfactor and MemQuality
// TrainingEfactor is a median between oldEfactor and ComputeEfactor
func (mem *Mem) ComputeTrainingEfactor(oldEfactor float32, quality MemQuality) {
	mem.ComputeEfactor(oldEfactor, quality)
	computedTrainingEfactor := (oldEfactor + mem.Efactor) / 2
	if computedTrainingEfactor < 1.3 {
		mem.Efactor = 1.3
	} else {
		mem.Efactor = computedTrainingEfactor
	}
}

// GetCardType returns the current CardType
// The CardType is CardMCQ if internal conditions are matched.
// Otherwise, it's Card.Type
func (mem *Mem) GetCardType() CardType {
	if mem.IsMCQ() {
		return CardMCQ
	}

	return mem.Card.Type
}

// ComputeInterval calculates and sets the interval between reviews
func (mem *Mem) ComputeInterval(oldInterval uint, eFactor float32, repetition uint) {
	switch repetition {
	case 0:
		mem.Interval = 1
	case 1, 2:
		mem.Interval = 2
	case 3:
		mem.Interval = 3
	default:
		mem.Interval = uint(float32(oldInterval)*eFactor*0.75) + 1
	}
}

func (mem *Mem) ComputeLearningStage() {
	switch {
	case mem.Repetition > 3 && mem.Repetition < 7:
		mem.LearningStage = StageReviewing
	case mem.Repetition > 7:
		mem.LearningStage = StageKnown
	default:
		mem.LearningStage = StageLearning
	}
}

// ComputeQualitySuccess sets the answer Quality
func (mem *Mem) ComputeQualitySuccess() {
	switch {
	case mem.GetCardType() == CardMCQ || mem.LearningStage == StageToLearn:
		mem.Quality = MemQualityError
	case mem.LearningStage == StageKnown:
		mem.Quality = MemQualityPerfect
	default:
		mem.Quality = MemQualityGoodMCQ
	}
}

// ComputeQualityFail sets the answer Quality
func (mem *Mem) ComputeQualityFail() {
	switch {
	case mem.GetCardType() == CardMCQ:
		if mem.LearningStage == StageToLearn {
			mem.Quality = MemQualityBlackout
		} else {
			mem.Quality = MemQualityErrorMCQ
		}
	case mem.LearningStage == StageLearning:
		mem.Quality = MemQualityErrorMCQ
	default:
		mem.Quality = MemQualityErrorHints
	}
}

// IsMCQ returns if the Mem should be an MCQ or not.
// It doesn't include Card.Type checks
func (mem *Mem) IsMCQ() bool {
	return mem.LearningStage < StageReviewing || mem.Efactor <= 1.7 || mem.Repetition < 2 || (mem.Efactor <= 2.3 && mem.Repetition < 3)
}
