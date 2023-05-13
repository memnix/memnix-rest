package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/random"
	"gorm.io/gorm"
)

// Deck is the domain model for a deck
type Deck struct {
	gorm.Model  `swaggerignore:"true"` // ignore this field when generating swagger docs
	Name        string                 `json:"name"`                           // Name of the deck
	Description string                 `json:"description"`                    // Description of the deck
	Lang        string                 `json:"lang"`                           // Lang of the deck
	Key         string                 `json:"key"`                            // Key of the deck
	Banner      string                 `json:"banner"`                         // Banner of the deck
	Learners    []*User                `json:"-" gorm:"many2many:user_decks;"` // Learners of the deck
	OwnerID     uint                   `json:"owner_id"`                       // OwnerID of the deck
	Status      DeckStatus             `json:"status"`                         // Status of the deck
}

// DeckStatus is the status of the deck
type DeckStatus int64

const (
	DeckStatusPrivate  DeckStatus = iota // DeckStatusPrivate is the private status of the deck
	DeckStatusToReview                   // DeckStatusToReview is the to review status of the deck
	DeckStatusPublic                     // DeckStatusPublic is the public status of the deck
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
	validate := validator.New()
	return validate.Struct(c)
}

// ToDeck converts the CreateDeck struct to a Deck struct
func (c *CreateDeck) ToDeck() Deck {
	key, _ := random.GenerateSecretCode(config.DeckSecretCodeLength)
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
