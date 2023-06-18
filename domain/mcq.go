package domain

import (
	"gorm.io/gorm"
	"strings"
)

// Mcq is the domain model for a mcq
type Mcq struct {
	gorm.Model `swaggerignore:"true"`
	Answers    string `json:"answers"`
	Linked     bool   `json:"linked"`
}

// IsLinked returns true if the mcq is linked.
func (m *Mcq) IsLinked() bool {
	return m.Linked
}

func (m *Mcq) ExtractAnswers() []string {
	// separate the answers by the semicolon
	return strings.Split(m.Answers, ";")
}
