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

type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type RegisterStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUser struct {
	Username    string     `json:"user_name"`
	Permissions Permission `json:"user_permissions" example:"0"` // 0: User; 1: Mod; 2: Admin
	Avatar      string     `json:"user_avatar" example:"avatar url"`
	Bio         string     `json:"user_bio" example:"A simple demo bio"`
}

func (publicUser *PublicUser) Set(user *User) {
	publicUser.Username = user.Username
	publicUser.Permissions = user.Permissions
	publicUser.Avatar = user.Avatar
	publicUser.Bio = user.Bio
}

// Permission enum type
type Permission int64

const (
	PermUser Permission = iota
	PermMod
	PermAdmin
)

// ToString returns Permission value as a string
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
