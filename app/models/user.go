package models

import (
	"gorm.io/gorm"
)

// User structure
type User struct {
	gorm.Model  `swaggerignore:"true"`
	Username    string     `json:"user_name" example:"Yume"`     // This should be unique
	Permissions Permission `json:"user_permissions" example:"0"` // 0: User; 1: Mod; 2: Admin
	Avatar      string     `json:"user_avatar" example:"avatar url"`
	Bio         string     `json:"user_bio" example:"A simple demo bio"`
	Email       string     `json:"email" gorm:"unique"`
	Password    []byte     `json:"-" swaggerignore:"true"`
}

type Permission int64

const (
	PermUser Permission = iota
	PermMod
	PermAdmin
)

func (s Permission) ToString() string {
	switch s {
	case PermUser:
		return "PermUser"
	case PermMod:
		return "PermMod"
	case PermAdmin:
		return "PermAdmin"
	default:
		return "Unknown"
	}
}
