package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GetAllCards method
// @Description Get every cards. Shouldn't really be used
// @Summary get all cards
// @Tags Card
// @Produce json
// @Success 200 {array} models.Card
// @Router /v1/cards/ [get]
func GetAllCards(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	var cards []models.Card

	if res := db.Joins("Deck").Find(&cards); res.Error != nil {

		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all cards",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All cards",
		Data:    cards,
		Count:   len(cards),
	})

}

// GetCardByID method to get a card by id
// @Description Get a card by tech id
// @Summary get a card
// @Tags Card
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} models.Card
// @Router /v1/cards/id/{id} [get]
func GetCardByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	// Params
	id := c.Params("id")

	card := new(models.Card)

	if err := db.Joins("Deck").First(&card, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get card by ID.",
		Data:    *card,
		Count:   1,
	})
}

// GetCardsFromDeck method to get cards from deck
// @Description Get every cards from a deck
// @Summary get a list of card
// @Tags Card
// @Produce json
// @Param deckID path int true "Deck ID"
// @Success 200 {array} models.Card
// @Router /v1/cards/deck/{deckID} [get]
func GetCardsFromDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("deckID")

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get cards from deck.",
		Data:    cards,
		Count:   len(cards),
	})
}

// POST

// CreateNewCard
func CreateNewCard(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	card := new(models.Card)

	if err := c.BodyParser(&card); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Preload("Deck").Create(card)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a card",
		Data:    *card,
		Count:   1,
	})
}
