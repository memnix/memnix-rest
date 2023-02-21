package domain

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `swaggerignore:"true"`
	Username   string     `json:"username"`                             // Username of the user
	Email      string     `json:"email" validate:"email" gorm:"unique"` // Email of the user
	Password   string     `json:"-"`                                    // Password of the user
	Avatar     string     `json:"avatar"`                               // Avatar of the user
	Permission Permission `json:"permission"`                           // Permission of the user
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
	ID         uint       `json:"id"`         // ID of the user
	Username   string     `json:"username"`   // Username of the user
	Email      string     `json:"email"`      // Email of the user
	Avatar     string     `json:"avatar"`     // Avatar of the user
	Permission Permission `json:"permission"` // Permission of the user
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
