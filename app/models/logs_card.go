package models

import (
	"gorm.io/gorm"
)

// Mem structure
type CardLogs struct {
	gorm.Model
	CardID uint `json:"card_id" example:"1"`
	Card   Card
	LogID  uint `json:"log_id" example:"1"`
	Log    Logs
}
