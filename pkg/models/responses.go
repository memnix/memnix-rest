package models

import (
	"gorm.io/gorm"
)

// Answer structure
type Answer struct {
	gorm.Model
	CardID uint `json:"card_id" example:"1"`
	Card   Card
	Answer string `json:"answer" example:"42"`
}

// DeckResponse structure
type DeckResponse struct {
	DeckID uint           `json:"deck_id"`
	Deck   Deck           `json:"deck"`
	Cards  []ResponseCard `json:"cards"`
	Count  int            `json:"count"`
}

type TodayResponse struct {
	DecksReponses []DeckResponse `json:"decks_responses"`
	Count         int            `json:"count"`
}

func (today *TodayResponse) AppendDeckResponse(deckResponse DeckResponse) {
	today.DecksReponses = append(today.DecksReponses, deckResponse)
}
