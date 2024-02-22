package domain_test

import (
	"testing"

	"github.com/memnix/memnix-rest/domain"
)

func TestCard_String(t *testing.T) {
	card := &domain.Card{
		Question: "What is the capital of France?",
		Answer:   "Paris",
	}

	expected := "What is the capital of France? Paris"
	result := card.String()

	if result != expected {
		t.Errorf("Expected result to be %q, but got %q", expected, result)
	}
}
func TestCard_TableName(t *testing.T) {
	card := &domain.Card{}
	expected := "cards"
	result := card.TableName()

	if result != expected {
		t.Errorf("Expected result to be %q, but got %q", expected, result)
	}
}

func TestCard_IsBlankType(t *testing.T) {
	testCases := []struct {
		name     string
		cardType domain.CardType
		expected bool
	}{
		{
			name:     "CardTypeBlankOnly",
			cardType: domain.CardTypeBlankOnly,
			expected: true,
		},
		{
			name:     "CardTypeBlankProgressive",
			cardType: domain.CardTypeBlankProgressive,
			expected: true,
		},
		{
			name:     "CardTypeBlankMCQ",
			cardType: domain.CardTypeBlankMCQ,
			expected: true,
		},
		{
			name:     "CardTypeQAOnly",
			cardType: domain.CardTypeQAOnly,
			expected: false,
		},
		{
			name:     "CardTypeMCQOnly",
			cardType: domain.CardTypeMCQOnly,
			expected: false,
		},
		{
			name:     "CardTypeQAProgressive",
			cardType: domain.CardTypeQAProgressive,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card := &domain.Card{
				CardType: tc.cardType,
			}

			result := card.IsBlankType()

			if result != tc.expected {
				t.Errorf("Expected result to be %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestCard_HasMcqType(t *testing.T) {
	testCases := []struct {
		name     string
		cardType domain.CardType
		expected bool
	}{
		{
			name:     "CardTypeBlankOnly",
			cardType: domain.CardTypeBlankOnly,
			expected: false,
		},
		{
			name:     "CardTypeBlankProgressive",
			cardType: domain.CardTypeBlankProgressive,
			expected: true,
		},
		{
			name:     "CardTypeBlankMCQ",
			cardType: domain.CardTypeBlankMCQ,
			expected: true,
		},
		{
			name:     "CardTypeQAOnly",
			cardType: domain.CardTypeQAOnly,
			expected: false,
		},
		{
			name:     "CardTypeMCQOnly",
			cardType: domain.CardTypeMCQOnly,
			expected: true,
		},
		{
			name:     "CardTypeQAProgressive",
			cardType: domain.CardTypeQAProgressive,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card := &domain.Card{
				CardType: tc.cardType,
			}

			result := card.HasMcqType()

			if result != tc.expected {
				t.Errorf("Expected result to be %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestCard_IsMcqType(t *testing.T) {
	testCases := []struct {
		name     string
		cardType domain.CardType
		expected bool
	}{
		{
			name:     "CardTypeBlankOnly",
			cardType: domain.CardTypeBlankOnly,
			expected: false,
		},
		{
			name:     "CardTypeBlankProgressive",
			cardType: domain.CardTypeBlankProgressive,
			expected: false,
		},
		{
			name:     "CardTypeBlankMCQ",
			cardType: domain.CardTypeBlankMCQ,
			expected: true,
		},
		{
			name:     "CardTypeQAOnly",
			cardType: domain.CardTypeQAOnly,
			expected: false,
		},
		{
			name:     "CardTypeMCQOnly",
			cardType: domain.CardTypeMCQOnly,
			expected: true,
		},
		{
			name:     "CardTypeQAProgressive",
			cardType: domain.CardTypeQAProgressive,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card := &domain.Card{
				CardType: tc.cardType,
			}

			result := card.IsMcqType()

			if result != tc.expected {
				t.Errorf("Expected result to be %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestCard_IsQAOnlyType(t *testing.T) {
	testCases := []struct {
		name     string
		cardType domain.CardType
		expected bool
	}{
		{
			name:     "CardTypeBlankOnly",
			cardType: domain.CardTypeBlankOnly,
			expected: false,
		},
		{
			name:     "CardTypeBlankProgressive",
			cardType: domain.CardTypeBlankProgressive,
			expected: false,
		},
		{
			name:     "CardTypeBlankMCQ",
			cardType: domain.CardTypeBlankMCQ,
			expected: false,
		},
		{
			name:     "CardTypeQAOnly",
			cardType: domain.CardTypeQAOnly,
			expected: true,
		},
		{
			name:     "CardTypeMCQOnly",
			cardType: domain.CardTypeMCQOnly,
			expected: false,
		},
		{
			name:     "CardTypeQAProgressive",
			cardType: domain.CardTypeQAProgressive,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card := &domain.Card{
				CardType: tc.cardType,
			}

			result := card.IsQAOnlyType()

			if result != tc.expected {
				t.Errorf("Expected result to be %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestCardType_cardTypeNames(t *testing.T) {
	expected := map[domain.CardType]string{
		domain.CardTypeQAOnly:           "QAOnly",
		domain.CardTypeMCQOnly:          "MCQOnly",
		domain.CardTypeQAProgressive:    "QAProgressive",
		domain.CardTypeBlankOnly:        "BlankOnly",
		domain.CardTypeBlankProgressive: "BlankProgressive",
		domain.CardTypeBlankMCQ:         "BlankMCQ",
	}

	for k, v := range expected {
		result := k.String()
		if result != v {
			t.Errorf("Expected result to be %q, but got %q", v, result)
		}
	}
}
