package models

import (
	"gorm.io/gorm"
)

// Mem structure
type UserLogs struct {
	gorm.Model
	UserID  uint `json:"user_id" example:"1"`
	User    User
	LogType UserLogType `json:"user_log_type" example:"1"`
	Message string      `json:"user_message" example:"Edited Profile Picture"`
}

type UserLogType int64

const (
	UserLogin UserLogType = iota
	UserLogout
	UserRegister
	UserEdit
	UserDeleted
	UserSubscribe
	UserUnsubscribe
)
