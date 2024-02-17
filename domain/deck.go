package domain

import (
	"github.com/memnix/memnix-rest/pkg/random"
	"gorm.io/gorm"
)

const DeckSecretCodeLength = 10

// Deck is the domain model for a deck
type Deck struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Lang        string     `json:"lang"`
	Key         string     `json:"key"`
	Banner      string     `json:"banner"`
	Learners    []*User    `json:"-" gorm:"many2many:user_decks;"`
	Cards       []Card     `json:"cards" gorm:"foreignKey:DeckID"`
	OwnerID     uint       `json:"owner_id"`
	Status      DeckStatus `json:"status"`
}

// TableName returns the table name for the deck model
func (*Deck) TableName() string {
	return "decks"
}

// DeckStatus is the status of the deck
type DeckStatus int64

const (
	DeckStatusPrivate  DeckStatus = 0 // DeckStatusPrivate is the private status of the deck
	DeckStatusToReview DeckStatus = 1 // DeckStatusToReview is the to review status of the deck
	DeckStatusPublic   DeckStatus = 2 // DeckStatusPublic is the public status of the deck
)

// PublicDeck is the public deck model
type PublicDeck struct {
	Name        string `json:"name"`        // Name of the deck
	Description string `json:"description"` // Description of the deck
	Lang        string `json:"lang"`        // Lang of the deck
	Banner      string `json:"banner"`      // Banner of the deck
	ID          uint   `json:"id"`          // ID of the deck
}

// ToPublicDeck converts the deck to a public deck
func (d *Deck) ToPublicDeck() PublicDeck {
	return PublicDeck{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Lang:        d.Lang,
		Banner:      d.Banner,
	}
}

// IsOwner checks if the deck is owned by the user with the given id
func (d *Deck) IsOwner(id uint) bool {
	return d.OwnerID == id
}

// CreateDeck is a struct that contains the data needed to create a deck
type CreateDeck struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Lang        string `json:"lang" validate:"required,len=2,alpha"`
	Banner      string `json:"banner" validate:"required,uri"`
}

// Validate validates the CreateDeck struct
func (c *CreateDeck) Validate() error {
	return validate.Struct(c)
}

// ToDeck converts the CreateDeck struct to a Deck struct
func (c *CreateDeck) ToDeck() Deck {
	key, _ := random.GenerateSecretCode(DeckSecretCodeLength)
	return Deck{
		Name:        c.Name,
		Description: c.Description,
		Lang:        c.Lang,
		Status:      DeckStatusPrivate,
		Key:         key,
	}
}

// DeckIndex is the index of the deck for MeiliSearch search engine
type DeckIndex map[string]interface{}
