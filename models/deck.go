package models

import (
	"gorm.io/gorm"
)

// User structure
type Deck struct {
	gorm.Model
	DeckName string `json:"deck_name" example:"First Deck"`
}
