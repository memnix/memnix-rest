package domain

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model    `swaggerignore:"true"`
	Username      string     `json:"username"`
	Email         string     `json:"email" validate:"email" gorm:"unique"`
	Password      string     `json:"-"`
	Avatar        string     `json:"avatar"`
	OauthProvider string     `json:"oauth_provider" `
	OauthID       string     `json:"oauth_id" gorm:"unique"`
	Permission    Permission `json:"permission"`
	Oauth         bool       `json:"oauth" gorm:"default:false"`
}

func (u *User) ToPublicUser() PublicUser {
	return PublicUser{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Avatar:     u.Avatar,
		Permission: u.Permission,
	}
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) HasPermission(permission Permission) bool {
	return u.Permission >= permission
}

type PublicUser struct {
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Avatar     string     `json:"avatar"`
	ID         uint       `json:"id"`
	Permission Permission `json:"permission"`
}

type Login struct {
	Email    string `json:"email" validate:"email"` // Email of the user
	Password string `json:"password"`               // Password of the user
}

type Register struct {
	Username string `json:"username" validate:"required"` // Username of the user
	Email    string `json:"email" validate:"email"`       // Email of the user
	Password string `json:"password" validate:"required"` // Password of the user
}

func (r *Register) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *Register) ToUser() User {
	return User{
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
		Permission: PermissionUser,
	}
}
