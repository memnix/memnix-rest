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
	Deck       Deck
}

type McqType int64

const (
	McqStandalone McqType = iota
	McqLinked
)

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

func (mcq *Mcq) ExtractAnswers() []string {
	return strings.Split(mcq.Answers, ";")
}

func (mcq *Mcq) SetAnswers(answers []string) {
	mcq.Answers = strings.Join(answers, ";")
}

func (mcq *Mcq) QueryLinkedAnswers() []string {
	db := database.DBConn // DB Conn
	var cards []Card
	var responses []string

	if err := db.Joins("Mcq").Where("cards.mcq_id = ?", mcq.ID).Find(&cards).Error; err != nil {
		return responses
		//TODO: Error logging
	}

	for i := range cards {
		responses = append(responses, cards[i].Answer)
	}

	return responses
}

func (mcq *Mcq) NotValidate() bool {

	return mcq.Type == McqStandalone && (len(mcq.Answers) < utils.MinMcqAnswersLen || len(mcq.Answers) > utils.MaxMcqAnswersLen) || len(mcq.Name) > utils.MaxMcqName || mcq.Name == ""
}

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
