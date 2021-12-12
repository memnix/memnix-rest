package models

import (
	"gorm.io/gorm"
)

// Access structure
type Rating struct {
	gorm.Model
	UserID uint `json:"user_id" example:"1"`
	User   User
	DeckID uint `json:"deck_id" example:"1"`
	Deck   Deck
	Value  uint `json:"value" example:"3" gorm:"check:value < 6"`
}
