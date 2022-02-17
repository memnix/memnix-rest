package models

import (
	"gorm.io/gorm"
)

// Logs structure
type Logs struct {
	gorm.Model
	LogType LogType `json:"log_type" example:"1"`
	Message string  `json:"log_message" example:"Edited Profile Picture"`
}

type LogType int64

const (
	LogUndefined LogType = iota
	LogUserLogin
	LogUserLogout
	LogUserRegister
	LogUserEdit
	LogUserDeleted
	LogSubscribe
	LogUnsubscribe
	LogDeckCreated
	LogDeckDeleted
	LogDeckEdited
	LogDeckRated
	LogCardCreated
	LogCardDeleted
	LogCardEdited
)

// CardLogs structure
type CardLogs struct {
	gorm.Model
	CardID uint `json:"card_id" example:"1"`
	Card   Card
	LogID  uint `json:"log_id" example:"1"`
	Log    Logs
}

// DeckLogs structure
type DeckLogs struct {
	gorm.Model
	DeckID uint `json:"deck_id" example:"1"`
	Deck   Deck
	LogID  uint `json:"log_id" example:"1"`
	Log    Logs
}

// UserLogs structure
type UserLogs struct {
	gorm.Model
	UserID uint `json:"user_id" example:"1"`
	User   User
	LogID  uint `json:"log_id" example:"1"`
	Log    Logs
}
