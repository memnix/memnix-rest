package models

import (
	"time"

	"gorm.io/gorm"
)

// MemDate structure
type MemDate struct {
	gorm.Model    `swaggerignore:"true"`
	UserID        uint          `json:"user_id" example:"1"`
	User          User          `swaggerignore:"true"`
	CardID        uint          `json:"card_id" example:"1"`
	Card          Card          `swaggerignore:"true"`
	DeckID        uint          `json:"deck_id" example:"1"`
	Deck          Deck          `swaggerignore:"true"`
	NextDate      time.Time     `json:"next_date" example:"01/01/2000"` // gorm:"autoCreateTime"`
	LearningStage LearningStage `json:"learning_stage" gorm:"default:0"`
}

// LearningStage enum type
type LearningStage int64

const (
	StageNeverSeen LearningStage = iota
	StageToLearn
	StageToRelearn
	StageLearning
	StageReviewing
	StageKnown
)

// ComputeNextDate calculates and sets the NextDate
func (m *MemDate) ComputeNextDate(interval int) {
	m.NextDate = time.Now().AddDate(0, 0, interval)
}

// SetDefaultNextDate fills MemDate values and sets NextDate as time.Now()
func (m *MemDate) SetDefaultNextDate(userID, cardID, deckID uint) {
	m.UserID = userID
	m.CardID = cardID
	m.DeckID = deckID
	m.NextDate = time.Now()
}
