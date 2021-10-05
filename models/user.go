
package models

import (
	"gorm.io/gorm"
)

// User structure
type User struct {
	gorm.Model
	UserName string `json:"user_name" example:"Yume"`
}
