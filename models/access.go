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
	Permission uint `json:"permission" example:"0"` // 0: Student - 1: Contributor - 2: Editor - 3: Owner
}
