package controllers

import (
	"github.com/gofiber/fiber/v2"
	"memnixrest/app/models"
	queries2 "memnixrest/app/queries"
	"memnixrest/pkg/database"
	"memnixrest/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

// GetTodayCard method
// @Description Get next today card
// @Summary gets a card
// @Tags Card
// @Produce json
// @Success 200 {object} models.Card
// @Router /v1/cards/today [get]
func GetTodayCard(c *fiber.Ctx) error {
	res := new(models.ResponseHTTP)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	if res = queries2.FetchNextCard(auth.User.ID, 0, true); !res.Success {
		return queries2.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get today's card",
		Data:    res.Data,
		Count:   1,
	})
}

// GetTrainingCardsByDeck method
// @Description Get training cards
// @Summary gets a list of cards
// @Tags Card
// @Produce json
// @Success 200 {array} models.Card
// @Router /v1/cards/{deckID}/training [get]
func GetTrainingCardsByDeck(c *fiber.Ctx) error {
	res := new(models.ResponseHTTP)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	deckID := c.Params("deckID")
	deckIdInt, _ := strconv.ParseInt(deckID, 10, 32)

	access := queries2.CheckAccess(auth.User.ID, uint(deckIdInt), models.AccessStudent)
	if !access.Success {
		return queries2.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if res = queries2.FetchTrainingCards(auth.User.ID, uint(deckIdInt)); !res.Success {
		return queries2.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get today's card",
		Data:    res,
		Count:   1,
	})
}

// GetNextCard method
// @Description Get next card
// @Summary gets a card
// @Tags Card
// @Produce json
// @Success 200 {object} models.Card
// @Router /v1/cards/next [get]
func GetNextCard(c *fiber.Ctx) error {

	res := new(models.ResponseHTTP)
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	if res = queries2.FetchNextCard(auth.User.ID, 0, false); !res.Success {
		return queries2.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get next card",
		Data:    res.Data,
		Count:   1,
	})
}

// GetNextCardByDeck method
// @Description Get next card by deckID
// @Summary get a card
// @Tags Card
// @Produce json
// @Success 200 {object} models.Card
// @Router /v1/cards/{deckID}/next [get]
func GetNextCardByDeck(c *fiber.Ctx) error {

	deckID := c.Params("deckID")
	deckIDInt, _ := strconv.Atoi(deckID)
	res := new(models.ResponseHTTP)
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	if res = queries2.FetchNextCard(auth.User.ID, uint(deckIDInt), false); !res.Success {
		return queries2.RequestError(c, http.StatusInternalServerError, res.Message)

	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get next card by deck",
		Data:    res.Data,
		Count:   1,
	})
}

// GetAllCards method
// @Description Get every card. Shouldn't really be used
// @Summary gets all cards
// @Tags Card
// @Produce json
// @Success 200 {array} models.Card
// @Router /v1/cards/ [get]
func GetAllCards(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	var cards []models.Card

	if res := db.Joins("Deck").Find(&cards); res.Error != nil {
		return queries2.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
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
// @Summary gets a card
// @Tags Card
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} models.Card
// @Router /v1/cards/id/{id} [get]
func GetCardByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}
	// Params
	id := c.Params("id")

	card := new(models.Card)

	if err := db.Joins("Deck").First(&card, id).Error; err != nil {
		return queries2.RequestError(c, http.StatusInternalServerError, err.Error())

	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get card by ID.",
		Data:    *card,
		Count:   1,
	})
}

// GetCardsFromDeck method to get cards from deck
// @Description Get every card from a deck
// @Summary gets a list of card
// @Tags Card
// @Produce json
// @Param deckID path int true "Deck ID"
// @Success 200 {array} models.Card
// @Router /v1/cards/deck/{deckID} [get]
func GetCardsFromDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("deckID")

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
		return queries2.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get cards from deck.",
		Data:    cards,
		Count:   len(cards),
	})
}

// POST

// CreateNewCard method
// @Description Create a new card
// @Summary creates a card
// @Tags Card
// @Produce json
// @Accept json
// @Param card body models.Card true "Card to create"
// @Success 200
// @Router /v1/cards/new [post]
func CreateNewCard(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	result := new(models.ResponseHTTP)
	card := new(models.Card)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	if err := c.BodyParser(&card); err != nil {
		return queries2.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries2.CheckAccess(auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		return queries2.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if len(card.Question) <= 5 || len(card.Answer) <= 5 {
		return queries2.RequestError(c, http.StatusBadRequest, utils.ErrorQALen)
	}

	db.Create(card)

	log := queries2.CreateLog(models.LogCardCreated, auth.User.Username+" created "+card.Question)
	_ = queries2.CreateUserLog(auth.User.ID, log)
	_ = queries2.CreateDeckLog(card.DeckID, log)
	_ = queries2.CreateCardLog(card.ID, log)

	var users []models.User

	if result = queries2.GetSubUsers(card.DeckID); !result.Success {
		return queries2.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}

	switch _ := result.Data.(type) {
	default:
		return queries2.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	case []models.User:
		users = result.Data.([]models.User)
	}

	for _, s := range users {
		_ = queries2.GenerateMemDate(s.ID, card.ID, card.DeckID)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a card",
		Data:    *card,
		Count:   1,
	})
}

// PostResponse method
// @Description Post a response and check it
// @Summary posts a response
// @Tags Card
// @Produce json
// @Success 200
// @Accept json
// @Router /v1/cards/response [post]
func PostResponse(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	response := new(models.CardResponse)
	card := new(models.Card)

	if err := c.BodyParser(&response); err != nil {
		return queries2.RequestError(c, http.StatusBadRequest, err.Error())

	}

	if err := db.Joins("Deck").First(&card, response.CardID).Error; err != nil {
		return queries2.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	res := queries2.CheckAccess(auth.User.ID, card.Deck.ID, models.AccessStudent)
	if !res.Success {
		return queries2.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	validation := new(models.CardResponseValidation)

	if strings.EqualFold(
		strings.ReplaceAll(response.Response, " ", ""), strings.ReplaceAll(card.Answer, " ", "")) {
		validation.Validate = true
		validation.Message = "Correct answer"
	} else {
		validation.Validate = false
		validation.Message = "Incorrect answer"
	}

	_ = queries2.PostMem(&auth.User, card, validation, response.Training)

	validation.Answer = card.Answer

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success post response",
		Data:    *validation,
		Count:   1,
	})
}

// PUT

// UpdateCardByID method
// @Description Edit a card
// @Summary edits a card
// @Tags Card
// @Produce json
// @Success 200
// @Accept json
// @Param card body models.Card true "card to edit"
// @Router /v1/cards/{cardID}/edit [put]
func UpdateCardByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return queries2.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := queries2.CheckAccess(auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		return queries2.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if err := UpdateCard(c, card); !err.Success {
		return queries2.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := queries2.CreateLog(models.LogCardEdited, auth.User.Username+" edited "+card.Question)
	_ = queries2.CreateUserLog(auth.User.ID, log)
	_ = queries2.CreateCardLog(card.ID, log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update card by ID",
		Data:    *card,
		Count:   1,
	})
}

// UpdateCard function
func UpdateCard(c *fiber.Ctx, card *models.Card) *models.ResponseHTTP {
	db := database.DBConn

	deckId := card.DeckID

	res := new(models.ResponseHTTP)

	if err := c.BodyParser(&card); err != nil {
		res.GenerateError(err.Error())
		return res
	}

	if deckId != card.DeckID {
		res.GenerateError(utils.ErrorBreak)
		return res
	}

	if len(card.Question) <= 5 || len(card.Answer) <= 5 {
		res.GenerateError(utils.ErrorQALen)
		return res
	}

	db.Save(card)

	res.GenerateSuccess("Success update card", nil, 0)
	return res
}

// DeleteCardById method
// @Description Delete a card
// @Summary deletes a card
// @Tags Card
// @Produce json
// @Success 200
// @Router /v1/cards/{cardID} [delete]
func DeleteCardById(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	id := c.Params("id")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries2.AuthError(c, auth)
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return queries2.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	if res := queries2.CheckAccess(auth.User.ID, card.DeckID, models.AccessOwner); !res.Success {
		return queries2.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	db.Delete(card)

	log := queries2.CreateLog(models.LogCardDeleted, auth.User.Username+" deleted "+card.Question)
	_ = queries2.CreateUserLog(auth.User.ID, log)
	_ = queries2.CreateCardLog(card.ID, log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success delete card by ID",
		Data:    *card,
		Count:   1,
	})

}
