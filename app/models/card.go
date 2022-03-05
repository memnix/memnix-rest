package models

import (
	"github.com/memnix/memnixrest/pkg/database"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// Card structure
type Card struct {
	gorm.Model  `swaggerignore:"true"`
	Question    string `json:"card_question" example:"What's the answer to life ?"`
	Answer      string `json:"card_answer" example:"42"`
	DeckID      uint   `json:"deck_id" example:"1"`
	Deck        Deck
	Tips        string   `json:"card_tips" example:"The answer is from a book"`
	Explication string   `json:"card_explication" example:"The number 42 is the answer to life has written in a very famous book"`
	Type        CardType `json:"card_type" example:"0" gorm:"type:Int"`
	Format      string   `json:"card_format" example:"Date / Name / Country"`
	Image       string   `json:"card_image"` // Should be an url
}

// CardType enum type
type CardType int64

const (
	CardString CardType = iota
	CardInt
	CardMCQ
)

// ToString returns CardType value as a string
func (s CardType) ToString() string {
	switch s {
	case CardString:
		return "Card String"
	case CardInt:
		return "Card Int"
	case CardMCQ:
		return "Card MCQ"
	default:
		return "Unknown"
	}
}

// GetMCQAnswers returns 3 random incorrect MCQAnswers
func (card *Card) GetMCQAnswers() []string {
	db := database.DBConn // DB Conn
	var answersList []string
	var answers []Answer

	if err := db.Joins("Card").Where("answers.card_id = ?", card.ID).Limit(3).Order("random()").Find(&answers).Error; err != nil {
		return nil
	}

	if len(answers) >= 3 {
		answersList = append(answersList, answers[0].Answer, answers[1].Answer, answers[2].Answer, card.Answer)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(answersList), func(i, j int) { answersList[i], answersList[j] = answersList[j], answersList[i] })
	}

	return answersList
}
