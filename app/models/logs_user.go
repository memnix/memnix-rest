package models

import (
	"gorm.io/gorm"
)

// UserLogs structure
type UserLogs struct {
	gorm.Model
	UserID uint `json:"user_id" example:"1"`
	User   User
	LogID  uint `json:"log_id" example:"1"`
	Log    Logs
}
