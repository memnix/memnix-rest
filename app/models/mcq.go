package models

import (
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

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
			mcq.UpdateLinkedAnswers()
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
	db := database.DBConn // DB Conn
	var cards []Card

	if err := db.Joins("Mcq").Where("cards.mcq_id = ?", mcq.ID).Find(&cards).Error; err != nil {
		return make([]string, 0)
		//TODO: Error logging
	}

	responses := make([]string, len(cards))
	for i := range cards {
		responses = append(responses, cards[i].Answer)
	}

	return responses
}

// NotValidate performs validation of the mcq
func (mcq *Mcq) NotValidate() bool {
	return mcq.Type == McqStandalone && (len(mcq.Answers) < utils.MinMcqAnswersLen || len(mcq.Answers) > utils.MaxMcqAnswersLen) || len(mcq.Name) > utils.MaxMcqName || mcq.Name == ""
}

// FillWithLinkedAnswers method
func (mcq *Mcq) FillWithLinkedAnswers() *ResponseHTTP {
	res := new(ResponseHTTP)

	answers := mcq.QueryLinkedAnswers()
	if len(answers) == 0 {
		res.GenerateError("Couldn't query linked answers")
		return res
	}
	mcq.SetAnswers(answers)

	res.GenerateSuccess("Success fill mcq with linked answers", answers, len(answers))
	return res
}

// UpdateLinkedAnswers method to update the db
func (mcq *Mcq) UpdateLinkedAnswers() *ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(ResponseHTTP)

	if err := mcq.FillWithLinkedAnswers(); !err.Success {
		res.GenerateError(err.Message)
		return res
	}

	if err := db.Save(mcq).Error; err != nil {
		res.GenerateError(err.Error())
		return res
	}

	res.GenerateSuccess("Success update mcq with linked answers", nil, 0)
	return res
}
