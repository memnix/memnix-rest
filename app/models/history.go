package models

import (
	"gorm.io/gorm"
)

// History structure
type History struct {
	gorm.Model
	UserID      uint `json:"user_id" example:"1"`
	User        User
	CardID      uint `json:"card_id" example:"1"`
	Card        Card
	Description string `json:"description" example:"Rewrite question with more details"`
}
