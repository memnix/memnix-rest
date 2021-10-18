package models

import (
	"gorm.io/gorm"
)

// Deck structure
type Deck struct {
	gorm.Model
	DeckName string `json:"deck_name" example:"First Deck"`
	Status   uint   `json:"status" example:"0"` // 0: Draft - 1: Private - 2: Published
}
