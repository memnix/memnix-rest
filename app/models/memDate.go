package models

import (
	"memnixrest/pkg/database"
	"time"

	"gorm.io/gorm"
)

// MemDate structure
type MemDate struct {
	gorm.Model
	UserID   uint `json:"user_id" example:"1"`
	User     User
	CardID   uint `json:"card_id" example:"1"`
	Card     Card
	DeckID   uint `json:"deck_id" example:"1"`
	Deck     Deck
	NextDate time.Time `json:"next_date" example:"06/01/2003"` // gorm:"autoCreateTime"`
}

func (m *MemDate) ComputeNextDate(interval int) {
	m.NextDate = time.Now().AddDate(0, 0, interval)
}

func (m *MemDate) Generate(userID uint, cardID uint, deckID uint) {
	m.UserID = userID
	m.CardID = cardID
	m.DeckID = deckID
	m.NextDate = time.Now()
}

func (m *MemDate) GetNextToday(userID uint) *ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(ResponseHTTP)

	// Get next card with date condition
	t := time.Now()
	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		userID).Joins("Card").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.next_date < ? AND accesses.permission >= ?",
		userID, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour), AccessStudent).Limit(1).Order("next_date asc").Find(&m).Error; err != nil {
		res.GenerateError("Next today memDate not found")
		return res
	}

	res.GenerateSuccess("Success getting next today memDate", nil, 0)

	return res
}

func (m *MemDate) GetNext(userID uint) *ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(ResponseHTTP)

	// Get next card
	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		userID).Joins("Card").Joins("Deck").Where("mem_dates.user_id = ? AND accesses.permission >= ?",
		userID, AccessStudent).Limit(1).Order("next_date asc").Find(&m).Error; err != nil {
		res.GenerateError("Next memDate not found")
		return res
	}

	res.GenerateSuccess("Success getting next memDate", nil, 0)

	return res
}

func (m *MemDate) GetNextByDeck(userID uint, deckID uint) *ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(ResponseHTTP)

	// Get next card  with deck condition
	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.deck_id = ?",
		userID, deckID).Limit(1).Order("next_date asc").Find(&m).Error; err != nil {
		res.GenerateError("Next memDate by deck not found")
		return res
	}

	res.GenerateSuccess("Success getting next memDate by deck", nil, 0)

	return res
}
