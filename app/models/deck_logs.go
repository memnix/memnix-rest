package models

import (
	"gorm.io/gorm"
)

// Mem structure
type DeckLogs struct {
	gorm.Model
	DeckID  uint `json:"deck_id" example:"1"`
	Deck    Deck
	LogType DeckLogType `json:"deck_log_type" example:"1"`
	Message string      `json:"deck_message" example:"Added a new card"`
}

type DeckLogType int64

const (
	DeckCreated DeckLogType = iota
	DeckDeleted
	DeckEdited
)
