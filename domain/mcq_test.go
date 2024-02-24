package domain_test

import (
	"reflect"
	"testing"

	"github.com/memnix/memnix-rest/domain"
)

func TestMcq_TableName(t *testing.T) {
	mcq := &domain.Mcq{}

	expected := "mcqs"
	tableName := mcq.TableName()

	if tableName != expected {
		t.Errorf("Expected table name to be %s, but got %s", expected, tableName)
	}
}

func TestMcq_IsLinked(t *testing.T) {
	mcq := &domain.Mcq{
		Linked: true,
	}

	expected := true
	isLinked := mcq.IsLinked()

	if isLinked != expected {
		t.Errorf("Expected IsLinked to be %v, but got %v", expected, isLinked)
	}

	mcq2 := &domain.Mcq{
		Linked: false,
	}

	expected2 := false
	isLinked2 := mcq2.IsLinked()

	if isLinked2 != expected2 {
		t.Errorf("Expected IsLinked to be %v, but got %v", expected2, isLinked2)
	}
}

func TestMcq_ExtractAnswers(t *testing.T) {
	testCases := []struct {
		name     string
		answers  string
		expected []string
	}{
		{
			name:     "SingleAnswer",
			answers:  "A",
			expected: []string{"A"},
		},
		{
			name:     "MultipleAnswers",
			answers:  "A;B;C",
			expected: []string{"A", "B", "C"},
		},
		{
			name:     "NoAnswers",
			answers:  "",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mcq := &domain.Mcq{
				Answers: tc.answers,
			}

			result := mcq.ExtractAnswers()

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected ExtractAnswers to return %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestMcq_AppendAnswer(t *testing.T) {
	mcq := &domain.Mcq{
		Answers: "A",
	}

	answer := "B"
	expected := "A;B"
	mcq.AppendAnswer(answer)

	if mcq.Answers != expected {
		t.Errorf("Expected Answers to be %s, but got %s", expected, mcq.Answers)
	}
}

func TestMcq_RemoveAnswer(t *testing.T) {
	testCases := []struct {
		name     string
		answers  string
		answer   string
		expected string
	}{
		{
			name:     "SingleAnswer",
			answers:  "A",
			answer:   "A",
			expected: "",
		},
		{
			name:     "MultipleAnswers",
			answers:  "A;B;C",
			answer:   "B",
			expected: "A;C",
		},
		{
			name:     "NoAnswers",
			answers:  "",
			answer:   "A",
			expected: "",
		},
		{
			name:     "AnswerNotInAnswers",
			answers:  "A;B;C",
			answer:   "D",
			expected: "A;B;C",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mcq := &domain.Mcq{
				Answers: tc.answers,
			}

			mcq.RemoveAnswer(tc.answer)

			if mcq.Answers != tc.expected {
				t.Errorf("Expected Answers to be %s, but got %s", tc.expected, mcq.Answers)
			}
		})
	}
}

func TestMcq_UpdateAnswer(t *testing.T) {
	testCases := []struct {
		name      string
		answers   string
		oldAnswer string
		newAnswer string
		expected  string
	}{
		{
			name:      "SingleAnswer",
			answers:   "A",
			oldAnswer: "A",
			newAnswer: "B",
			expected:  "B",
		},
		{
			name:      "MultipleAnswers",
			answers:   "A;B;C",
			oldAnswer: "B",
			newAnswer: "D",
			expected:  "A;D;C",
		},
		{
			name:      "NoAnswers",
			answers:   "",
			oldAnswer: "A",
			newAnswer: "B",
			expected:  "",
		},
		{
			name:      "AnswerNotInAnswers",
			answers:   "A;B;C",
			oldAnswer: "D",
			newAnswer: "E",
			expected:  "A;B;C",
		},
		{
			name:      "AnswerNotInAnswers",
			answers:   "A;B;C",
			oldAnswer: "B",
			newAnswer: "A",
			expected:  "A;A;C",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mcq := &domain.Mcq{
				Answers: tc.answers,
			}

			mcq.UpdateAnswer(tc.oldAnswer, tc.newAnswer)

			if mcq.Answers != tc.expected {
				t.Errorf("Expected Answers to be %s, but got %s", tc.expected, mcq.Answers)
			}
		})
	}
}
