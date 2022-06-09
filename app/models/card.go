package models

import (
	"database/sql"
	"fmt"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// Card structure
type Card struct {
	gorm.Model       `swaggerignore:"true"`
	Question         string        `json:"card_question" example:"What's the answer to life ?"`
	Answer           string        `json:"card_answer" example:"42"`
	DeckID           uint          `json:"deck_id" example:"1"`
	Deck             Deck          `swaggerignore:"true" json:"-"`
	Type             CardType      `json:"card_type" example:"0" gorm:"type:Int"`
	Format           string        `json:"card_format" example:"Date / Name / Country"`
	Image            string        `json:"card_image"` // Should be an url
	Case             bool          `json:"card_case" gorm:"default:false"`
	Spaces           bool          `json:"card_spaces" gorm:"default:false"`
	Explication      string        `json:"card_explication"`
	ExplicationImage string        `json:"card_explication_image"`
	McqID            sql.NullInt32 `json:"mcq_id" swaggerignore:"true"`
	Mcq              Mcq
}

// CardType enum type
type CardType int64

const (
	CardString CardType = iota
	CardInt
	CardMCQ
)

// NotValidate performs validation of the Card
func (card *Card) NotValidate() bool {
	return len(card.Question) < utils.MinCardQuestionLen || card.Answer == "" || (card.Type == CardMCQ && card.McqID.Int32 == 0) || len(
		card.Format) > utils.MaxCardFormatLen || len(
		card.Question) > utils.MaxDefaultLen || len(card.Answer) > utils.MaxDefaultLen || len(card.Image) > utils.MaxImageUrlLen || len(
		card.ExplicationImage) > utils.MaxImageUrlLen || len(card.Explication) > utils.MaxCardExplicationLen
}

// ValidateMCQ makes sure that the mcq attached to a card is correct
func (card *Card) ValidateMCQ(user *User) (*Mcq, bool) {
	db := database.DBConn // DB Conn
	mcq := new(Mcq)

	if card.McqID.Int32 != 0 {
		if err := db.First(&mcq, card.McqID).Error; err != nil {
			log := CreateLog(fmt.Sprintf("Error on CreateNewCard: %s from %s", err.Error(), user.Email), LogQueryGetError).SetType(LogTypeError).AttachIDs(user.ID, card.DeckID, 0)
			_ = log.SendLog()
			return nil, false
		}
		if mcq.DeckID != card.DeckID {
			log := CreateLog(fmt.Sprintf("Error on CreateNewCard: card.DeckID != mcq.DeckID from %s", user.Email), LogBadRequest).SetType(LogTypeError).AttachIDs(user.ID, card.DeckID, 0)
			_ = log.SendLog()
			return nil, false
		}

		if mcq.Type == McqLinked {
			return mcq, true
		}
	}
	return nil, true
}

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

// GetMCQAnswers returns mcq's answers
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

	if len(answersList) < 3 {
		return answers
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(answersList), func(i, j int) { answersList[i], answersList[j] = answersList[j], answersList[i] })

	i, c := 0, 0
	for i < 3 {
		if answersList[c] != card.Answer {
			answers = append(answers, answersList[c])
			i++
		}
		c++
	}

	answers = append(answers, card.Answer)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	return answers
}
