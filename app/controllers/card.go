package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/app/queries"
	"github.com/memnix/memnixrest/pkg/core"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"net/http"
	"strconv"
)

// GetTodayCard method
// @Description Get next today card
// @Summary gets a card
// @Tags Card
// @Produce json
// @Success 200 {object} models.Card
// @Deprecated
// @Router /v1/cards/today/one [get]
func GetTodayCard(c *fiber.Ctx) error {
	res := new(models.ResponseHTTP)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	if res = queries.FetchNextTodayCard(auth.User.ID); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Error from %s on GetTodayCard: %s", auth.User.Email, res.Message), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get today's card",
		Data:    res.Data,
		Count:   1,
	})
}

// GetAllTodayCard method
// @Description Get all today card
// @Summary gets a list of card
// @Tags Card
// @Produce json
// @Success 200 {array} models.Card
// @Router /v1/cards/today [get]
func GetAllTodayCard(c *fiber.Ctx) error {
	res := new(models.ResponseHTTP)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	if res = queries.FetchTodayCard(auth.User.ID); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Error on GetAllTodayCard: %s", res.Message), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
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
		return queries.AuthError(c, &auth)
	}

	deckID := c.Params("deckID")
	deckIdInt, _ := strconv.ParseInt(deckID, 10, 32)

	access := queries.CheckAccess(auth.User.ID, uint(deckIdInt), models.AccessStudent)
	if !access.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - GetTodayCard: %s", auth.User.Email, deckIdInt, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, uint(deckIdInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if res = queries.FetchTrainingCards(auth.User.ID, uint(deckIdInt)); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Error on GetTrainingCardsByDeck: %s from %s", res.Message, auth.User.Email), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, uint(deckIdInt), 0)
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

// GetNextCard method
// @Description Get next card
// @Summary gets a card
// @Tags Card
// @Deprecated
// @Produce json
// @Success 200 {object} models.Card
// @Router /v1/cards/next [get]
func GetNextCard(c *fiber.Ctx) error {

	res := new(models.ResponseHTTP)
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	if res = queries.FetchNextCard(auth.User.ID, 0); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Error on GetNextCard: %s from %s", res.Message, auth.User.Email), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)
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
// @Deprecated
// @Success 200 {object} models.Card
// @Router /v1/cards/{deckID}/next [get]
func GetNextCardByDeck(c *fiber.Ctx) error {

	deckID := c.Params("deckID")
	deckIDInt, _ := strconv.Atoi(deckID)
	res := new(models.ResponseHTTP)
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	if res = queries.FetchNextCard(auth.User.ID, uint(deckIDInt)); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Error on GetNextCardByDeck: %s from %s on deck %d", res.Message, auth.User.Email, deckIDInt), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, uint(deckIDInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)

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
		return queries.AuthError(c, &auth)
	}

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
		return queries.AuthError(c, &auth)
	}
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
// @Success 200 {array} models.Card
// @Router /v1/cards/deck/{deckID} [get]
func GetCardsFromDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("deckID")
	deckID, _ := strconv.ParseUint(id, 10, 32)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	if res := queries.CheckAccess(auth.User.ID, uint(deckID), models.AccessEditor); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - GetCardsFromDeck: %s", auth.User.Email, deckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, uint(deckID), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	var cards []models.Card

	if err := db.Joins("Deck").Joins("Mcq").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on GetCardsFromDeck: %s from %s on %d", err.Error(), auth.User.Email, deckID), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, uint(deckID), 0)
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
	card := new(models.Card)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	if err := c.BodyParser(&card); err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on CreateNewCard: %s from %s", err.Error(), auth.User.Email), models.LogBodyParserError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries.CheckAccess(auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateNewCard: %s", auth.User.Email, card.DeckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if res := queries.CheckCardLimit(auth.User.Permissions, card.DeckID); !res {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateNewCard: This deck has reached his limit", auth.User.Email, card.DeckID), models.LogDeckCardLimit).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, "This deck has reached his limit ! You can't add more card to it.")
	}

	if card.NotValidate() {
		log := models.CreateLog(fmt.Sprintf("BadRequest from %s on deck %d - CreateNewCard: BadRequest", auth.User.Email, card.DeckID), models.LogBadRequest).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorQALen)
	}

	_, ok := card.ValidateMCQ(&auth.User)
	if !ok {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}

	db.Create(card)

	log := models.CreateLog(fmt.Sprintf("Created: %d - %s", card.ID, card.Question), models.LogCardCreated).SetType(models.LogTypeInfo).AttachIDs(auth.User.ID, card.DeckID, card.ID)
	_ = log.SendLog()

	if res := queries.UpdateSubUsers(card, &auth.User); res != nil {
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
// @Success 200
// @Accept json
// @Router /v1/cards/selfresponse [post]
func PostSelfEvaluateResponse(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	response := new(models.CardSelfResponse)
	card := new(models.Card)

	if err := c.BodyParser(&response); err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on PostSelfEvaluateResponse: %s from %s", err.Error(), auth.User.Email), models.LogBodyParserError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if response.Quality > 4 || response.Quality < 1 {
		log := models.CreateLog(fmt.Sprintf("Error on PostSelfEvaluateResponse: Quality > 4 from %s", auth.User.Email), models.LogBodyParserError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorBreak)
	}

	if err := db.Joins("Deck").First(&card, response.CardID).Error; err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on PostSelfEvaluateResponse: %s from %s", err.Error(), auth.User.Email), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	res := queries.CheckAccess(auth.User.ID, card.Deck.ID, models.AccessStudent)
	if !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - PostSelfEvaluateResponse: %s", auth.User.Email, card.DeckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, card.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	//TODO: Add error handling
	_ = queries.PostSelfEvaluatedMem(&auth.User, card, response.Quality, response.Training)

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
// @Success 200
// @Accept json
// @Router /v1/cards/response [post]
func PostResponse(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	response := new(models.CardResponse)
	card := new(models.Card)

	if err := c.BodyParser(&response); err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on PostResponse: %s from %s", err.Error(), auth.User.Email), models.LogBodyParserError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if err := db.Joins("Deck").First(&card, response.CardID).Error; err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on PostResponse: %s from %s", err.Error(), auth.User.Email), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	res := queries.CheckAccess(auth.User.ID, card.Deck.ID, models.AccessStudent)
	if !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - PostResponse: %s", auth.User.Email, card.DeckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, card.DeckID, 0)
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
	_ = queries.PostMem(&auth.User, card, validation, response.Training)

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
	cardID, _ := strconv.ParseUint(id, 10, 32)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {

		log := models.CreateLog(fmt.Sprintf("Error on UpdateCardByID: %s from %s", err.Error(), auth.User.Email), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := queries.CheckAccess(auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - UpdateCardByID: %s", auth.User.Email, card.DeckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, card.DeckID, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if err := UpdateCard(c, card, &auth.User); !err.Success {
		log := models.CreateLog(fmt.Sprintf("Error on UpdateCardByID: %s from %s", err.Message, auth.User.Email), models.LogBadRequest).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := models.CreateLog(fmt.Sprintf("Edited: %d - %s", card.ID, card.Question), models.LogCardEdited).SetType(models.LogTypeInfo).AttachIDs(auth.User.ID, card.DeckID, card.ID)
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

	if card.NotValidate() {
		res.GenerateError(utils.ErrorQALen)
		return res
	}

	shouldUpdateMcq := false
	mcq := new(models.Mcq)

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
	cardID, _ := strconv.ParseUint(id, 10, 32)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	if res := queries.CheckAccess(auth.User.ID, card.DeckID, models.AccessOwner); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - DeleteCardById: %s", auth.User.Email, card.DeckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, card.DeckID, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	var memDates []models.MemDate

	if err := db.Joins("Card").Where("mem_dates.card_id = ?", card.ID).Find(&memDates).Error; err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on DeleteCardById: %s from %s", err.Error(), auth.User.Email), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, uint(cardID))
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
		// TODO: Error
	}

	db.Unscoped().Delete(memDates)

	db.Delete(card)

	log := models.CreateLog(fmt.Sprintf("Deleted: %d - %s", card.ID, card.Question), models.LogCardDeleted).SetType(models.LogTypeInfo).AttachIDs(auth.User.ID, card.DeckID, card.ID)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success delete card by ID",
		Data:    *card,
		Count:   1,
	})

}
