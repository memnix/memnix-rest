package domain

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`               // Username of the user
	Email    string `json:"email" validate:"email"` // Email of the user
	Password string `json:"password"`               // Password of the user
}

func (u *User) Validate() error {
	var validate = validator.New()
	return validate.Struct(u)
}
