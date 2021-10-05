package models

import (
	"gorm.io/gorm"
)

// User structure
type Revision struct {
	gorm.Model
	UserID uint `json:"user_id" example:"1"`
	User   User
	CardID uint `json:"card_id" example:"1"`
	Card   Card
	Result bool `json:"result" example:"true"` // True means that the answer was given right.

}
