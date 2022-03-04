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
	Quality    uint    `json:"quality" example:"0"` // [0: Blackout - 1: Error with choices - 2: Error with hints - 3: Error - 4: Good with hints - 5: Perfect]
	Repetition uint    `json:"repetition" example:"0" `
	Efactor    float32 `json:"e_factor" example:"2.5"`
	Interval   uint    `json:"interval" example:"0"`
}

func (m *Mem) ComputeEfactor(oldEfactor float32, quality uint) {
	eFactor := oldEfactor + (0.1 - (5.0-float32(quality))*(0.08+(5-float32(quality)))*0.02)

	if eFactor < 1.3 {
		m.Efactor = 1.3
	} else {
		m.Efactor = eFactor
	}
}

func (m *Mem) ComputeTrainingEfactor(oldEfactor float32, quality uint) {
	eFactor := oldEfactor + (0.1 - (5.0-float32(quality))*(0.08+(5-float32(quality)))*0.02)
	computedEfactor := (oldEfactor + eFactor) / 2
	if computedEfactor < 1.3 {
		m.Efactor = 1.3
	} else {
		m.Efactor = computedEfactor
	}
}

func (m *Mem) GetMemType() CardType {
	if m.Efactor <= 2 || m.Repetition < 2 || (m.Efactor <= 2.3 && m.Repetition < 4) {
		return CardMCQ
	}

	return m.Card.Type
}

func (m *Mem) ComputeInterval(oldInterval uint, eFactor float32, repetition uint) {
	switch repetition {
	case 0:
		m.Interval = 1
	case 1, 2:
		m.Interval = 2
	case 3:
		m.Interval = 3
	default:
		m.Interval = uint(float32(oldInterval)*eFactor*0.75) + 1
	}
}
