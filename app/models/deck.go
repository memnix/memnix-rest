package models

import (
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

// Deck structure
type Deck struct {
	gorm.Model  `swaggerignore:"true"`
	DeckName    string     `json:"deck_name" example:"First Deck"`
	Description string     `json:"deck_description" example:"A simple demo deck"`
	Banner      string     `json:"deck_banner" example:"A banner url"`
	Status      DeckStatus `json:"deck_status" example:"2"` // 1: Draft - 2: Private - 3: Published
	Key         string     `json:"deck_key" example:"MEM"`
	Code        string     `json:"deck_code" example:"6452"`
	Share       bool       `json:"deck_share" example:"true" gorm:"default:false"`
}

// DeckStatus enum type
type DeckStatus int64

const (
	DeckPrivate DeckStatus = iota + 1
	DeckWaitingReview
	DeckPublic
)

// ToString returns DeckStatus value as a string
func (s DeckStatus) ToString() string {
	switch s {
	case DeckWaitingReview:
		return "Deck Waiting Review"
	case DeckPrivate:
		return "Deck Private"
	case DeckPublic:
		return "Deck Public"
	default:
		return "Unknown"
	}
}

// NotValidate performs validation of the deck
func (deck *Deck) NotValidate() bool {
	return len(deck.DeckName) < utils.MinDeckLen || len(deck.DeckName) > utils.MaxDeckNameLen || len(deck.Description) < utils.MinDeckLen || len(deck.Description) > utils.MaxDefaultLen || len(deck.Banner) > utils.MaxImageUrlLen // || len(strings.ReplaceAll(deck.DeckName, " ", "")) != utils.DeckKeyLen
}

// GenerateCode creates a random code from the deck key
func (deck *Deck) GenerateCode() {
	rand.Seed(time.Now().UTC().UnixNano())

	randomInt := rand.Intn(99)
	runes := []rune(deck.Key)
	var result []int
	result = append(result, randomInt)

	for i := 0; i < len(runes); i++ {
		result = append(result, int(runes[i]))
	}

	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })

	deck.Code = strconv.Itoa(int(result[0]+result[1])/2) + strconv.Itoa(result[4]) + strconv.Itoa(int(result[2]+result[3])/2)
}

// GetOwner returns the deck Owner
func (deck *Deck) GetOwner() User {
	db := database.DBConn

	access := new(Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.deck_id =? AND accesses.permission >= ?", deck.ID, AccessOwner).Find(&access).Error; err != nil {
		return access.User
	}

	return access.User
}
