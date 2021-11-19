package models

import (
	"time"

	"gorm.io/gorm"
)

// Mem structure
type MemDate struct {
	gorm.Model
	UserID   uint `json:"user_id" example:"1"`
	User     User
	CardID   uint `json:"card_id" example:"1"`
	Card     Card
	DeckID   uint `json:"deck_id" example:"1"`
	Deck     Deck
	NextDate time.Time `json:"next_date" example:"06/01/2003"` // gorm:"autoCreateTime"`
}
