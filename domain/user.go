package domain

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `swaggerignore:"true"`
	Username   string `json:"username"`                             // Username of the user
	Email      string `json:"email" validate:"email" gorm:"unique"` // Email of the user
	Password   string `json:"-"`                                    // Password of the user
	Avatar     string `json:"avatar" validate:"avatar"`             // Avatar of the user
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
