package models

import (
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/utils"
	"gorm.io/gorm"
	"strings"
)

// Answer structure
type Answer struct {
	gorm.Model
	CardID uint `json:"card_id" example:"1"`
	Card   Card
	Answer string `json:"answer" example:"42"`
}

type Mcq struct {
	gorm.Model `swaggerignore:"true"`
	Name       string  `json:"mcq_name"`
	Answers    string  `json:"mcq_answers"`
	Type       McqType `json:"mcq_type"`
	DeckID     uint    `json:"deck_id" example:"1"`
	Deck       Deck    `swaggerignore:"true" json:"-"`
}

type McqType int64

const (
	McqStandalone McqType = iota
	McqLinked
)

// GetAnswers returns a list of answers
func (mcq *Mcq) GetAnswers() []string {
	answersList := mcq.ExtractAnswers()

	if mcq.Type == McqLinked {
		answers := mcq.QueryLinkedAnswers()
		if len(answers) != len(answersList) {
			//TODO: FIX THIS viewmodels.UpdateLinkedAnswers(mcq)
			return answers
		}
	}

	return answersList
}

// ExtractAnswers method
func (mcq *Mcq) ExtractAnswers() []string {
	return strings.Split(mcq.Answers, ";")
}

// SetAnswers method
func (mcq *Mcq) SetAnswers(answers []string) {
	mcq.Answers = strings.Join(answers, ";")
}

// QueryLinkedAnswers returns linked answers
func (mcq *Mcq) QueryLinkedAnswers() []string {
	db := infrastructures.GetDBConn() // DB Conn
	var cards []Card

	if err := db.Joins("Mcq").Where("cards.mcq_id = ?", mcq.ID).Find(&cards).Error; err != nil {
		return make([]string, 0)
		//TODO: Error logging
	}

	responses := make([]string, len(cards))
	for i := range cards {
		responses[i] = cards[i].Answer
	}

	return responses
}

// NotValidate performs validation of the mcq
func (mcq *Mcq) NotValidate() bool {
	return mcq.Type == McqStandalone && (len(mcq.Answers) < utils.MinMcqAnswersLen || len(mcq.Answers) > utils.MaxMcqAnswersLen) || len(mcq.Name) > utils.MaxMcqName || mcq.Name == ""
}
