package models

import (
	"gorm.io/gorm"
)

// Mem structure
type DeckLogs struct {
	gorm.Model
	DeckID uint `json:"deck_id" example:"1"`
	Deck   Deck
	LogID  uint `json:"log_id" example:"1"`
	Log    Logs
}
