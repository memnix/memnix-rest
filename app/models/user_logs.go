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
	UserLogin UserLogType = iota + 1
	UserLogout
	UserRegister
	UserEdit
	UserDeleted
	UserSubscribe
	UserUnsubscribe
)

func (s UserLogType) ToString() string {
	switch s {
	case UserLogin:
		return "User Login"
	case UserLogout:
		return "User Logout"
	case UserRegister:
		return "User Register"
	case UserEdit:
		return "User Edit"
	case UserDeleted:
		return "User Deleted"
	case UserSubscribe:
		return "User Subscribe"
	case UserUnsubscribe:
		return "User Unsubscribe"
	default:
		return "Unknown"
	}
}
