package models

import (
	"gorm.io/gorm"
)

// Card structure
type Card struct {
	gorm.Model
	Question    string `json:"card_question" example:"What's the answer to life ?"`
	Answer      string `json:"card_answer" example:"42"`
	DeckID      uint   `json:"deck_id" example:"1"`
	Deck        Deck
	Tips        string `json:"card_tips" example:"The answer is from a book"`
	Explication string `json:"card_explication" example:"The number 42 is the answer to life has written in a very famous book"`
}
