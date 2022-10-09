package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/pkg/core"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/logger"
	"github.com/memnix/memnixrest/pkg/models"
	"github.com/memnix/memnixrest/pkg/queries"
	"github.com/memnix/memnixrest/pkg/utils"
	"net/http"
	"strconv"
)

// GetAllTodayCard function to get all today card for a user
// @Description Get all today card
// @Summary gets a list of card
// @Tags Card
// @Produce json
// @Success 200  {array} models.TodayResponse
// @Security Beaver
// @Router /v1/cards/today [get]
func GetAllTodayCard(c *fiber.Ctx) error {
	var res *models.ResponseHTTP

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	if res = queries.FetchTodayCard(user.ID); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Error on GetAllTodayCard: %s", res.Message), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get today's cards",
		Data:    res.Data,
		Count:   res.Count,
	})
}

// GetTrainingCardsByDeck function to get training cards by deck
// @Description Get training cards from a deck
// @Summary gets a list of cards
// @Tags Card
// @Produce json
// @Success 200 {array} models.Card
// @Param deckId path int true "Deck ID"
// @Security Beaver
// @Router /v1/cards/{deckID}/training [get]
func GetTrainingCardsByDeck(c *fiber.Ctx) error {
	res := new(models.ResponseHTTP)

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	deckID := c.Params("deckID")
	deckIDInt, _ := strconv.ParseInt(deckID, 10, 32)

	access := queries.CheckAccess(user.ID, uint(deckIDInt), models.AccessStudent)
	if !access.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - GetTodayCard: %s", user.Email, deckIDInt, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, uint(deckIDInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if res = queries.FetchTrainingCards(user.ID, uint(deckIDInt)); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Error on GetTrainingCardsByDeck: %s from %s", res.Message, user.Email), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, uint(deckIDInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get today's card",
		Data:    res.Data,
		Count:   res.Count,
	})
}

// GetAllCards function to get all cards (deprecated)
// @Description Get every card. Shouldn't really be used
// @Summary gets all cards
// @Tags Card
// @Produce json
// @Security Admin
// @Success 200 {array} models.Card
// @Router /v1/cards/ [get]
// @Deprecated
func GetAllCards(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var cards []models.Card

	if res := db.Joins("Deck").Find(&cards); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All cards",
		Data:    cards,
		Count:   len(cards),
	})
}

// GetCardByID function to get a card by id
// @Description Get a card by id
// @Summary gets a card
// @Tags Card
// @Produce json
// @Param id path int true "Card ID"
// @Security Admin
// @Success 200 {object} models.Card
// @Router /v1/cards/id/{id} [get]
func GetCardByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	card := new(models.Card)

	if err := db.Joins("Deck").Joins("Mcq").First(&card, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
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
// @Security Beaver
// @Success 200 {array} models.Card
// @Router /v1/cards/deck/{deckID} [get]
func GetCardsFromDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("deckID")
	deckID, _ := strconv.ParseUint(id, 10, 32)

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	if res := queries.CheckAccess(user.ID, uint(deckID), models.AccessStudent); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - GetCardsFromDeck: %s", user.Email, deckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, uint(deckID), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	var cards []models.Card

	if err := db.Joins("Deck").Joins("Mcq").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on GetCardsFromDeck: %s from %s on %d", err.Error(), user.Email, deckID), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, uint(deckID), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
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
// @Description Create a new card (must be a deck editor)
// @Summary creates a card
// @Tags Card
// @Produce json
// @Accept json
// @Security Beaver
// @Param card body models.Card true "Card to create"
// @Success 200
// @Router /v1/cards/new [post]
func CreateNewCard(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	card := new(models.Card)

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	if err := c.BodyParser(&card); err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on CreateNewCard: %s from %s", err.Error(), user.Email), logger.LogBodyParserError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries.CheckAccess(user.ID, card.DeckID, models.AccessEditor); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateNewCard: %s", user.Email, card.DeckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if res := queries.CheckCardLimit(user.Permissions, card.DeckID); !res {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateNewCard: This deck has reached his limit", user.Email, card.DeckID), logger.LogDeckCardLimit).SetType(logger.LogTypeWarning).AttachIDs(user.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, "This deck has reached his limit ! You can't add more card to it.")
	}

	if card.NotValidate() {
		log := logger.CreateLog(fmt.Sprintf("BadRequest from %s on deck %d - CreateNewCard: BadRequest", user.Email, card.DeckID), logger.LogBadRequest).SetType(logger.LogTypeWarning).AttachIDs(user.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorQALen)
	}

	_, ok = card.ValidateMCQ(&user)
	if !ok {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}

	db.Create(card)

	log := logger.CreateLog(fmt.Sprintf("Created: %d - %s", card.ID, card.Question), logger.LogCardCreated).SetType(logger.LogTypeInfo).AttachIDs(user.ID, card.DeckID, card.ID)
	_ = log.SendLog()

	if res := queries.UpdateSubUsers(card, &user); res != nil {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a card",
		Data:    *card,
		Count:   1,
	})
}

// PostSelfEvaluateResponse method
// @Description Post a self evaluated response
// @Summary posts a response
// @Tags Card
// @Produce json
// @Security Beaver
// @Accept json
// @Param card body models.CardSelfResponse true "Self response"
// @Success 200
// @Router /v1/cards/selfresponse [post]
func PostSelfEvaluateResponse(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	response := new(models.CardSelfResponse)
	card := new(models.Card)

	if err := c.BodyParser(&response); err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on PostSelfEvaluateResponse: %s from %s", err.Error(), user.Email), logger.LogBodyParserError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if response.Quality > 4 || response.Quality < 1 {
		log := logger.CreateLog(fmt.Sprintf("Error on PostSelfEvaluateResponse: Quality > 4 from %s", user.Email), logger.LogBodyParserError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorBreak)
	}

	if err := db.Joins("Deck").First(&card, response.CardID).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on PostSelfEvaluateResponse: %s from %s", err.Error(), user.Email), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	res := queries.CheckAccess(user.ID, card.Deck.ID, models.AccessStudent)
	if !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - PostSelfEvaluateResponse: %s", user.Email, card.DeckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	//TODO: Add error handling
	_ = queries.PostSelfEvaluatedMem(&user, card, response.Quality, response.Training)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success post response",
		Data:    nil,
		Count:   1,
	})
}

// PostResponse method
// @Description Post a response and check it
// @Summary posts a response
// @Tags Card
// @Produce json
// @Security Beaver
// @Accept json
// @Param card body models.CardResponse true "Response"
// @Success 200 {object} models.CardResponseValidation
// @Router /v1/cards/response [post]
func PostResponse(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	response := new(models.CardResponse)
	card := new(models.Card)

	if err := c.BodyParser(&response); err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on PostResponse: %s from %s", err.Error(), user.Email), logger.LogBodyParserError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if err := db.Joins("Deck").First(&card, response.CardID).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on PostResponse: %s from %s", err.Error(), user.Email), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	res := queries.CheckAccess(user.ID, card.Deck.ID, models.AccessStudent)
	if !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - PostResponse: %s", user.Email, card.DeckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	validation := new(models.CardResponseValidation)

	if core.ValidateAnswer(response.Response, card) {
		validation.SetCorrect()
	} else {
		validation.SetIncorrect()
	}

	//TODO: Add error handling
	_ = queries.PostMem(&user, card, validation, response.Training)

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
// @Success 200 {object} models.Card
// @Security Beaver
// @Accept json
// @Param card body models.Card true "card to edit"
// @Param id path int true "card id"
// @Router /v1/cards/{cardID}/edit [put]
func UpdateCardByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")
	cardID, _ := strconv.ParseUint(id, 10, 32)

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}
	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on UpdateCardByID: %s from %s", err.Error(), user.Email), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := queries.CheckAccess(user.ID, card.DeckID, models.AccessEditor); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - UpdateCardByID: %s", user.Email, card.DeckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, card.DeckID, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if err := UpdateCard(c, card, &user); !err.Success {
		log := logger.CreateLog(fmt.Sprintf("Error on UpdateCardByID: %s from %s", err.Message, user.Email), logger.LogBadRequest).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := logger.CreateLog(fmt.Sprintf("Edited: %d - %s", card.ID, card.Question), logger.LogCardEdited).SetType(logger.LogTypeInfo).AttachIDs(user.ID, card.DeckID, card.ID)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update card by ID",
		Data:    *card,
		Count:   1,
	})
}

// UpdateCard function
func UpdateCard(c *fiber.Ctx, card *models.Card, user *models.User) *models.ResponseHTTP {
	db := database.DBConn

	deckID := card.DeckID

	res := new(models.ResponseHTTP)

	if err := c.BodyParser(&card); err != nil {
		res.GenerateError(err.Error())
		return res
	}

	if deckID != card.DeckID {
		res.GenerateError(utils.ErrorBreak)
		return res
	}

	if card.NotValidate() {
		res.GenerateError(utils.ErrorQALen)
		return res
	}

	shouldUpdateMcq := false

	mcq, ok := card.ValidateMCQ(user)
	if !ok {
		res.GenerateError(utils.ErrorRequestFailed)
		return res
	}

	shouldUpdateMcq = mcq != nil

	db.Save(card)

	if shouldUpdateMcq {
		mcq.UpdateLinkedAnswers()
	}

	res.GenerateSuccess("Success update card", nil, 0)
	return res
}

// DeleteCardByID method
// @Description Delete a card (must be a deck owner)
// @Summary deletes a card
// @Tags Card
// @Produce json
// @Security Beaver
// @Param id path int true "card id"
// @Success 200
// @Router /v1/cards/{cardID} [delete]
func DeleteCardByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	id := c.Params("id")
	cardID, _ := strconv.ParseUint(id, 10, 32)

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	if res := queries.CheckAccess(user.ID, card.DeckID, models.AccessOwner); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - DeleteCardById: %s", user.Email, card.DeckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, card.DeckID, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	var memDates []models.MemDate

	if err := db.Joins("Card").Where("mem_dates.card_id = ?", card.ID).Find(&memDates).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on DeleteCardById: %s from %s", err.Error(), user.Email), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
		// TODO: Error
	}

	db.Unscoped().Delete(memDates)

	db.Delete(card)

	log := logger.CreateLog(fmt.Sprintf("Deleted: %d - %s", card.ID, card.Question), logger.LogCardDeleted).SetType(logger.LogTypeInfo).AttachIDs(user.ID, card.DeckID, card.ID)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success delete card by ID",
		Data:    *card,
		Count:   1,
	})
}
