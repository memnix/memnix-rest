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

// CardsResponse structure
type CardsResponse struct {
	CardID uint `json:"card_id""`
	Card   Card
	Stage  LearningStage `json:"learning_stage"`
}

// DeckResponse structure
type DeckResponse struct {
	DeckID uint            `json:"deck_id"`
	Cards  []CardsResponse `json:"cards"`
	Count  int32           `json:"count"`
}

type DeckResponsePreview struct {
	DeckID        uint `json:"deck_id"`
	Deck          Deck
	ToLearnCount  int32 `json:"to_learn_count"`
	LearningCount int32 `json:"learning_count"'`
	ToReviewCount int32 `json:"to_review_count"`
	TotalCount    int32 `json:"total_count"`
}

type ResponsePreview struct {
	Decks []DeckResponsePreview `json:"decks_response"`
	Count int32                 `json:"count"`
}
