package models

import (
	"gorm.io/gorm"
)

// User structure
type User struct {
	gorm.Model
	Username    string `json:"user_name" example:"Yume" gorm:"unique"`                  // This should be unique
	DiscordID   string `json:"user_discord" example:"282233191916634113" gorm:"unique"` // This is unique
	Permissions int    `json:"user_permissions" example:"0"`                            // 0: User; 1: Mod; 2: Admin
	Avatar      string `json:"user_avatar" example:"avatar url"`
	Bio         string `json:"user_bio" example:"A simple demo bio"`
}
