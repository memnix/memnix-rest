package models

import (
	"github.com/memnix/memnixrest/pkg/utils"
	"gorm.io/gorm"
)

// Access structure
type Access struct {
	gorm.Model  `swaggerignore:"true"`
	UserID      uint             `json:"user_id" example:"1"`
	User        User             `swaggerignore:"true"`
	DeckID      uint             `json:"deck_id" example:"1"`
	Deck        Deck             `swaggerignore:"true"`
	Permission  AccessPermission `json:"permission" example:"0"` // 0: None - 1: Student - 2: Editor - 3: Owner
	ToggleToday bool             `json:"today" gorm:"default:true"`
}

// AccessPermission  enum type
type AccessPermission uint8

const (
	AccessNone AccessPermission = iota
	AccessStudent
	AccessEditor
	AccessOwner
)

// ToString returns AccessPermission value as a string
func (s AccessPermission) ToString() string {
	switch s {
	case AccessStudent:
		return "Access Student"
	case AccessEditor:
		return "Access Editor"
	case AccessOwner:
		return "Access Owner"
	default:
		return utils.UNKNOWN
	}
}

// Set Access values
func (access *Access) Set(userID, deckID uint, permission AccessPermission) {
	access.UserID = userID
	access.DeckID = deckID
	access.Permission = permission
}

// TODO add setter for AccessPermission
