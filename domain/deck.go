package domain

import (
	"github.com/go-playground/validator/v10"
	"github.com/memnix/memnix-rest/pkg/random"
	"gorm.io/gorm"
)

type Deck struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Lang        string     `json:"lang"`
	Key         string     `json:"key"`
	Banner      string     `json:"banner"`
	Learners    []*User    ` json:"-" gorm:"many2many:user_decks;"`
	OwnerID     uint       `json:"owner_id"`
	Status      DeckStatus `json:"status"`
}

type DeckStatus int64

const (
	DeckStatusPrivate DeckStatus = iota
	DeckStatusToReview
	DeckStatusPublic
)

type PublicDeck struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Lang        string `json:"lang"`
	Banner      string `json:"banner"`
	ID          uint   `json:"id"`
}

func (d *Deck) ToPublicDeck() PublicDeck {
	return PublicDeck{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Lang:        d.Lang,
		Banner:      d.Banner,
	}
}

func (d *Deck) IsOwner(id uint) bool {
	return d.OwnerID == id
}

type CreateDeck struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Lang        string `json:"lang" validate:"required,len=2,alpha"`
	Banner      string `json:"banner" validate:"required,uri"`
}

func (c *CreateDeck) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (c *CreateDeck) ToDeck() Deck {
	key, _ := random.GenerateSecretCode(10)
	return Deck{
		Name:        c.Name,
		Description: c.Description,
		Lang:        c.Lang,
		Status:      DeckStatusPrivate,
		Key:         key,
	}
}

type DeckIndex map[string]interface{}
