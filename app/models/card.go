package models

import (
	"gorm.io/gorm"
)

// Card structure
type Card struct {
	gorm.Model  `swaggerignore:"true"`
	Question    string `json:"card_question" example:"What's the answer to life ?"`
	Answer      string `json:"card_answer" example:"42"`
	DeckID      uint   `json:"deck_id" example:"1"`
	Deck        Deck
	Tips        string `json:"card_tips" example:"The answer is from a book"`
	Explication string `json:"card_explication" example:"The number 42 is the answer to life has written in a very famous book"`
	Type        uint   `json:"card_type" example:"0"` // 0: int - 1: string - 2: MCQ //TODO: Add Enum
	Format      string `json:"card_format" example:"Date / Name / Country"`
	Image       string `json:"card_image"` // Should be an url
}
