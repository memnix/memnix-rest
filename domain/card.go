package domain

import (
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model `swaggerignore:"true"`
	Question   string   `json:"question"`
	Answer     string   `json:"answer"`
	Mcq        Mcq      `json:"mcq" gorm:"foreignKey:McqID"`
	DeckID     uint     `json:"deck_id"`
	McqID      uint     `json:"mcq_id"`
	CardType   CardType `json:"card_type"`
}

func (c *Card) TableName() string {
	return "cards"
}

// String returns the string representation of the card.
func (c *Card) String() string {
	return c.Question + " " + c.Answer
}

// HasMcqType returns true if the card has a mcq type.
// It returns true for both mcq and progressive types.
func (c *Card) HasMcqType() bool {
	return c.IsMcqType() || c.IsProgressiveType()
}

// IsLinked returns true if the card is linked.
// It returns true for both mcq and progressive types.
func (c *Card) IsLinked() bool {
	return c.IsMcqType() && c.Mcq.IsLinked()
}

// IsMcqType returns true if the card is a mcq type.
// It returns true only for mcq and blank mcq types.
func (c *Card) IsMcqType() bool {
	return c.CardType == CardTypeMCQOnly || c.CardType == CardTypeBlankMCQ
}

// IsProgressiveType returns true if the card is a progressive type.
func (c *Card) IsProgressiveType() bool {
	return c.CardType == CardTypeQAProgressive || c.CardType == CardTypeBlankProgressive
}

// IsBlankType returns true if the card is a blank type.
func (c *Card) IsBlankType() bool {
	return c.CardType == CardTypeBlankOnly || c.CardType == CardTypeBlankProgressive || c.CardType == CardTypeBlankMCQ
}

// IsQAOnlyType returns true if the card is a qa only type.
func (c *Card) IsQAOnlyType() bool {
	return c.CardType == CardTypeQAOnly
}

// CardType is the type of the card
type CardType int32

const (
	CardTypeQAOnly           CardType = iota // CardTypeQAOnly is the qa only type of the card
	CardTypeMCQOnly                          // CardTypeMCQOnly is the mcq only type of the card
	CardTypeQAProgressive                    // CardTypeQAProgressive is the qa progressive type of the card
	CardTypeBlankOnly                        // CardTypeBlankOnly is the blank only type of the card
	CardTypeBlankProgressive                 // CardTypeBlankProgressive is the blank progressive type of the card
	CardTypeBlankMCQ                         // CardTypeBlankMCQ is the blank mcq type of the card
)

// String returns the string representation of the card type.
func (c CardType) String() string {
	return [...]string{"QAOnly", "MCQOnly", "QAProgressive", "BlankOnly", "BlankProgressive", "BlankMCQ"}[c]
}
