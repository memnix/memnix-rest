package models

import (
	"gorm.io/gorm"
)

// Deck structure
type Deck struct {
	gorm.Model
	DeckName string `json:"deck_name" example:"First Deck"`
}
