package domain

import (
	"gorm.io/gorm"
)

// User is the domain model for a user
type User struct {
	gorm.Model    `swaggerignore:"true"`
	Username      string     `json:"username"`
	Email         string     `json:"email" validate:"email" gorm:"unique"`
	Password      string     `json:"-"`
	Avatar        string     `json:"avatar"`
	OauthProvider string     `json:"oauth_provider" `
	OauthID       string     `json:"oauth_id" gorm:"unique"`
	Learning      []*Deck    `json:"learning" gorm:"many2many:user_decks;"`
	OwnDecks      []Deck     `json:"own_decks" gorm:"foreignKey:OwnerID"`
	Permission    Permission `json:"permission"`
	Oauth         bool       `json:"oauth" gorm:"default:false"`
}

// TableName returns the table name for the user model
func (*User) TableName() string {
	return "users"
}

// ToPublicUser converts the user to a public user
func (u *User) ToPublicUser() PublicUser {
	return PublicUser{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Avatar:     u.Avatar,
		Permission: u.Permission,
	}
}

// Validate validates the user
func (u *User) Validate() error {
	return validate.Struct(u)
}

// HasPermission checks if the user has the given permission
func (u *User) HasPermission(permission Permission) bool {
	return u.Permission >= permission
}

// PublicUser is the public user model
type PublicUser struct {
	Username   string     `json:"username"`   // Username of the user
	Email      string     `json:"email"`      // Email of the user
	Avatar     string     `json:"avatar"`     // Avatar of the user
	ID         uint       `json:"id"`         // ID of the user
	Permission Permission `json:"permission"` // Permission of the user
}

// Login is the login model
type Login struct {
	Email    string `json:"email" validate:"email"` // Email of the user
	Password string `json:"password"`               // Password of the user
}

// Register is the register model
type Register struct {
	Username string `json:"username" validate:"required"` // Username of the user
	Email    string `json:"email" validate:"email"`       // Email of the user
	Password string `json:"password" validate:"required"` // Password of the user
}

// Validate validates the register model
func (r *Register) Validate() error {
	return validate.Struct(r)
}

// ToUser converts the register model to a user
func (r *Register) ToUser() User {
	return User{
		Username:   r.Username,
		Email:      r.Email,
		Password:   r.Password,
		Permission: PermissionUser,
	}
}
