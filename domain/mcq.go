package domain

import (
	"strings"
)

// Mcq is the domain model for a mcq.
type Mcq struct {
	Answers string `json:"answers"`
	ID      uint   `json:"id"`
	Linked  bool   `json:"linked"`
}

// TableName returns the table name for the mcq model.
func (*Mcq) TableName() string {
	return "mcqs"
}

// IsLinked returns true if the mcq is linked.
func (m *Mcq) IsLinked() bool {
	return m.Linked
}

func (m *Mcq) ExtractAnswers() []string {
	if m.Answers == "" {
		return []string{}
	}
	// separate the answers by the semicolon
	return strings.Split(m.Answers, ";")
}

func (m *Mcq) AppendAnswer(answer string) {
	m.Answers += ";" + answer
}

func (m *Mcq) RemoveAnswer(answer string) {
	answers := m.ExtractAnswers()
	for i, a := range answers {
		if a == answer {
			answers = append(answers[:i], answers[i+1:]...)
			break
		}
	}
	m.Answers = strings.Join(answers, ";")
}

func (m *Mcq) UpdateAnswer(oldAnswer, newAnswer string) {
	if m.Answers == "" {
		return
	}

	answers := strings.Split(m.Answers, ";")
	for i := range answers {
		if answers[i] == oldAnswer {
			answers[i] = newAnswer
			break
		}
	}
	m.Answers = strings.Join(answers, ";")
}
