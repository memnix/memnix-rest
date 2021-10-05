package models

import (
	"gorm.io/gorm"
)

// User structure
type Identifier struct {
	gorm.Model
	UserID    uint `json:"user_id" example:"2"`
	User      User
	DiscordID string `json:"discord_id" example:"282233191916634113"`
}
