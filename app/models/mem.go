package models

import (
	"gorm.io/gorm"
)

// Mem structure
type Mem struct {
	gorm.Model
	UserID     uint `json:"user_id" example:"1"`
	User       User
	CardID     uint `json:"card_id" example:"1"`
	Card       Card
	Quality    MemQuality `json:"quality" example:"0"` // [0: Blackout - 1: Error with choices - 2: Error with hints - 3: Error - 4: Good with hints - 5: Perfect]
	Repetition uint       `json:"repetition" example:"0" `
	Efactor    float32    `json:"e_factor" example:"2.5"`
	Interval   uint       `json:"interval" example:"0"`
}

type MemQuality int64

const (
	MemQualityBlackout MemQuality = iota
	MemQualityErrorMCQ
	MemQualityErrorHints
	MemQualityError
	MemQualityGoodMCQ
	MemQualityPerfect
)

func (mem *Mem) FillDefaultValues(userID uint, cardID uint) {
	mem.UserID = userID
	mem.CardID = cardID
	mem.Quality = MemQualityBlackout
	mem.Repetition = 0
	mem.Efactor = 2.5
	mem.Interval = 0
}

func (mem *Mem) ComputeEfactor(oldEfactor float32, quality MemQuality) {
	eFactor := oldEfactor + (0.1 - (5.0-float32(quality))*(0.08+(5-float32(quality)))*0.02)

	if eFactor < 1.3 {
		mem.Efactor = 1.3
	} else {
		mem.Efactor = eFactor
	}
}

func (mem *Mem) ComputeTrainingEfactor(oldEfactor float32, quality MemQuality) {
	eFactor := oldEfactor + (0.1 - (5.0-float32(quality))*(0.08+(5-float32(quality)))*0.02)
	computedEfactor := (oldEfactor + eFactor) / 2
	if computedEfactor < 1.3 {
		mem.Efactor = 1.3
	} else {
		mem.Efactor = computedEfactor
	}
}

func (mem *Mem) GetCardType() CardType {
	if mem.Efactor <= 2 || mem.Repetition < 2 || (mem.Efactor <= 2.3 && mem.Repetition < 4) {
		return CardMCQ
	}

	return mem.Card.Type
}

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

func (mem *Mem) ComputeQualitySuccess() {
	if mem.GetCardType() == CardMCQ {
		mem.Quality = MemQualityError
	} else {
		if mem.Repetition > 3 {
			mem.Quality = MemQualityPerfect
		}
		mem.Quality = MemQualityGoodMCQ
	}
}

func (mem *Mem) ComputeQualityFail() {
	if mem.GetCardType() == CardMCQ {
		if mem.Repetition <= 3 {
			mem.Quality = MemQualityBlackout
		}
		mem.Quality = MemQualityErrorMCQ
	}
	if mem.Repetition <= 4 {
		mem.Quality = MemQualityErrorMCQ
	}
	mem.Quality = MemQualityErrorHints
}

func (mem *Mem) IsMCQ() bool {
	return mem.Efactor <= 1.4 || mem.Quality <= MemQualityErrorMCQ || mem.Repetition < 2
}
