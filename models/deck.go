package models

import (
	"gorm.io/gorm"
)

// Deck structure
type Deck struct {
	gorm.Model
	DeckName    string `json:"deck_name" example:"First Deck"`
	Description string `json:"deck_description" example:"A simple demo deck"`
	Banner      string `json:"deck_banner" example:"A banner url"`
	Status      uint   `json:"deck_status" example:"0"` // 0: Draft - 1: Private - 2: Published
}
