package models

import (
	"gorm.io/gorm"
)

// Deck structure
type Deck struct {
	gorm.Model  `swaggerignore:"true"`
	DeckName    string     `json:"deck_name" example:"First Deck"`
	Description string     `json:"deck_description" example:"A simple demo deck"`
	Banner      string     `json:"deck_banner" example:"A banner url"`
	Status      DeckStatus `json:"deck_status" example:"0"` // 1: Draft - 2: Private - 3: Published
}

type DeckStatus int64

const (
	DeckDraft DeckStatus = iota + 1
	DeckPrivate
	DeckPublic
)

func (s DeckStatus) String() string {
	switch s {
	case DeckDraft:
		return "Deck Draft"
	case DeckPrivate:
		return "Deck Private"
	case DeckPublic:
		return "Deck Public"
	default:
		return "Unknown"
	}
}
