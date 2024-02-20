package domain

import (
	"gorm.io/gorm"
)

// Card represents a card in the domain model.
// It includes fields for the question, answer, multiple choice question (mcq),
// deck ID, mcq ID, and card type.
type Card struct {
	gorm.Model `swaggerignore:"true"` // Model from gorm package
	Question   string                 `json:"question"`                    // The question on the card
	Answer     string                 `json:"answer"`                      // The answer to the question on the card
	Mcq        Mcq                    `json:"mcq" gorm:"foreignKey:McqID"` // The multiple choice question associated with the card
	DeckID     uint                   `json:"deck_id"`                     // The ID of the deck the card belongs to
	McqID      uint                   `json:"mcq_id"`                      // The ID of the multiple choice question associated with the card
	CardType   CardType               `json:"card_type"`                   // The type of the card
}

// TableName returns the name of the table in the database.
func (*Card) TableName() string {
	return "cards"
}

// String returns a string representation of the card.
func (c *Card) String() string {
	return c.Question + " " + c.Answer
}

// HasMcqType checks if the card is of a type that includes a multiple choice question.
func (c *Card) HasMcqType() bool {
	return c.IsMcqType() || c.IsProgressiveType()
}

// IsLinked checks if the card is linked to a multiple choice question.
func (c *Card) IsLinked() bool {
	return c.IsMcqType() && c.Mcq.IsLinked()
}

// IsMcqType checks if the card is of a type that is only a multiple choice question.
func (c *Card) IsMcqType() bool {
	return c.CardType == CardTypeMCQOnly || c.CardType == CardTypeBlankMCQ
}

// IsProgressiveType checks if the card is of a progressive type.
func (c *Card) IsProgressiveType() bool {
	return c.CardType == CardTypeQAProgressive || c.CardType == CardTypeBlankProgressive
}

// IsBlankType checks if the card is of a blank type.
func (c *Card) IsBlankType() bool {
	return c.CardType == CardTypeBlankOnly || c.CardType == CardTypeBlankProgressive || c.CardType == CardTypeBlankMCQ
}

// IsQAOnlyType checks if the card is of a type that is only a question and answer.
func (c *Card) IsQAOnlyType() bool {
	return c.CardType == CardTypeQAOnly
}

// CardType represents the type of card.
type CardType int32

// Constants representing the different types of cards.
const (
	CardTypeQAOnly           CardType = 0
	CardTypeMCQOnly          CardType = 1
	CardTypeQAProgressive    CardType = 2
	CardTypeBlankOnly        CardType = 3
	CardTypeBlankProgressive CardType = 4
	CardTypeBlankMCQ         CardType = 5
)

// String returns the string representation of a CardType.
func (c CardType) String() string {
	return c.cardTypeNames()[c]
}

// cardTypeNames returns a map of CardType values to their string representations.
func (CardType) cardTypeNames() map[CardType]string {
	return map[CardType]string{
		CardTypeQAOnly:           "QAOnly",
		CardTypeMCQOnly:          "MCQOnly",
		CardTypeQAProgressive:    "QAProgressive",
		CardTypeBlankOnly:        "BlankOnly",
		CardTypeBlankProgressive: "BlankProgressive",
		CardTypeBlankMCQ:         "BlankMCQ",
	}
}
