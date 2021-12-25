package models

import (
	"gorm.io/gorm"
)

// Mem structure
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
	LogCardCreated
	LogCardDeleted
	LogCardEdited
)
