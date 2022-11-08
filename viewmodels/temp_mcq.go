package viewmodels

import (
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/models"
)

// FillWithLinkedAnswers method
func FillWithLinkedAnswers(mcq *models.Mcq) *ResponseHTTP {
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
func UpdateLinkedAnswers(mcq *models.Mcq) *ResponseHTTP {
	db := infrastructures.GetDBConn() // DB Conn
	res := new(ResponseHTTP)

	if err := FillWithLinkedAnswers(mcq); !err.Success {
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
