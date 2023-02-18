package domain

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Barista struct {
	gorm.Model `swaggerignore:"true" faker:"-"`
	Name       string `json:"name" example:"John" faker:"first_name"  validate:"min=2,max=50"`
	Email      string `json:"email" example:"kafejo@kafejo.dev" faker:"email" validate:"email"`
}

func (b *Barista) Validate() error {
	var validate = validator.New()
	return validate.Struct(b)
}
