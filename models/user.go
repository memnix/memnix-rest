package models

import (
	"gorm.io/gorm"
)

// User structure
type User struct {
	gorm.Model
	UserName  string `json:"user_name" example:"Yume"`
	DiscordID string `json:"discord_id" example:"282233191916634113"`
}
