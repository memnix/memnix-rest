package models

import (
	"gorm.io/gorm"
)

// Access structure
type Access struct {
	gorm.Model
	UserID     uint `json:"user_id" example:"1"`
	User       User
	DeckID     uint `json:"deck_id" example:"1"`
	Deck       Deck
	Permission AccessPermission `json:"permission" example:"0"` // 0: None - 1: Student - 2: Contributor - 3: Editor - 4: Owner
}

type AccessPermission int64

const (
	AccessStudent AccessPermission = iota + 1
	AccessContributor
	AccessEditor
	AccessOwner
)

func (s AccessPermission) String() string {
	switch s {
	case AccessStudent:
		return "Access Student"
	case AccessContributor:
		return "Access Contributor"
	case AccessEditor:
		return "Access Editor"
	case AccessOwner:
		return "Access Owner"
	default:
		return "Unknown"
	}
}
