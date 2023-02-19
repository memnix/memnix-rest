package domain

import (
	"github.com/edgedb/edgedb-go"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID       edgedb.UUID `edgedb:"id"`
	Username string      `json:"username" edgedb:"username"`            // Username of the user
	Email    string      `json:"email" validate:"email" edgedb:"email"` // Email of the user
	Password string      `json:"password" edgedb:"password"`            // Password of the user
}

func (u *User) Validate() error {
	var validate = validator.New()
	return validate.Struct(u)
}
