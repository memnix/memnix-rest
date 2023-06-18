package domain

import (
	"strings"

	"gorm.io/gorm"
)

// Mcq is the domain model for a mcq
type Mcq struct {
	gorm.Model `swaggerignore:"true"`
	Answers    string `json:"answers"`
	Linked     bool   `json:"linked"`
}

// TableName returns the table name for the mcq model
func (m *Mcq) TableName() string {
	return "mcqs"
}

// IsLinked returns true if the mcq is linked.
func (m *Mcq) IsLinked() bool {
	return m.Linked
}

func (m *Mcq) ExtractAnswers() []string {
	// separate the answers by the semicolon
	return strings.Split(m.Answers, ";")
}
