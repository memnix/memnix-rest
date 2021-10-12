package models

import (
	"gorm.io/gorm"
)

// Access structure
type Access struct {
	gorm.Model
	UserID     uint `json:"user_id" example:"1"`
	User       User
	DeckID     uint `json:"deck_id" example:"1"`
	Deck       Deck
	Permission uint `json:"permission" example:"0"` // 0: None - 1: Student - 2: Contributor - 3: Editor - 4: Owner
}
