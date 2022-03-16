package models

import (
	"github.com/memnix/memnixrest/pkg/database"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// Card structure
type Card struct {
	gorm.Model `swaggerignore:"true"`
	Question   string `json:"card_question" example:"What's the answer to life ?"`
	Answer     string `json:"card_answer" example:"42"`
	DeckID     uint   `json:"deck_id" example:"1"`
	Deck       Deck
	Type       CardType `json:"card_type" example:"0" gorm:"type:Int"`
	Format     string   `json:"card_format" example:"Date / Name / Country"`
	Image      string   `json:"card_image"` // Should be an url
	McqID      uint     `json:"mcq_id"`
	Mcq        Mcq
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

func (card *Card) GetMCQAnswers() []string {
	db := database.DBConn // DB Conn

	var answers []string

	mcq := new(Mcq)

	if err := db.First(&mcq, card.McqID).Error; err != nil {
		return answers
	}

	answersList := mcq.GetAnswers()

	for i := range answersList {
		if i >= len(answersList) {
			break
		}
		if answersList[i] == card.Answer {
			answersList[i] = answersList[len(answersList)-1]
			answersList = answersList[:len(answersList)-1]
		}
	}

	if len(answersList) >= 3 {
		for i := 0; i < 3; i++ {
			index := rand.Intn(len(answersList) - 1)
			answers = append(answers, answersList[index])
			answersList[index] = answersList[len(answersList)-1]
			answersList = answersList[:len(answersList)-1]
		}
		answers = append(answers, card.Answer)
	}
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	return answers
}
