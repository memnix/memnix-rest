package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"memnixrest/pkg/queries"
	"memnixrest/pkg/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetTodayCard method
// @Description Get next today card
// @Summary gets a card
// @Tags Card
// @Produce json
// @Success 200 {object} models.Card
// @Router /v1/cards/today [get]
func GetTodayCard(c *fiber.Ctx) error {
	res := *new(models.ResponseHTTP)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	if res = queries.FetchNextTodayCard(c, &auth.User); !res.Success {
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get today's card",
		Data:    res.Data,
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
	//db := database.DBConn // DB Conn

	res := *new(models.ResponseHTTP)
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	if res = queries.FetchNextCard(c, &auth.User); !res.Success {
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
// @Success 200 {object} models.Card
// @Router /v1/cards/{deckID}/next [get]
func GetNextCardByDeck(c *fiber.Ctx) error {
	//db := database.DBConn // DB Conn

	deckID := c.Params("deckID")

	res := *new(models.ResponseHTTP)
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	if res = queries.FetchNextCardByDeck(c, &auth.User, deckID); !res.Success {
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
		return queries.AuthError(c, auth)
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
		return queries.AuthError(c, auth)
	}
	// Params
	id := c.Params("id")

	card := new(models.Card)

	if err := db.Joins("Deck").First(&card, id).Error; err != nil {
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

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
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
		return queries.AuthError(c, auth)
	}

	if err := c.BodyParser(&card); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries.CheckAccess(c, auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if len(card.Question) <= 5 || len(card.Answer) <= 5 {
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorQALen)
	}

	db.Create(card)

	log := queries.CreateLog(models.LogCardCreated, auth.User.Username+" created "+card.Question)
	_ = queries.CreateUserLog(auth.User.ID, log)
	_ = queries.CreateDeckLog(card.DeckID, log)
	_ = queries.CreateCardLog(card.ID, log)

	var users []models.User

	if users = queries.GetSubUsers(c, card.DeckID); len(users) > 0 {
		for _, s := range users {
			_ = queries.GenerateMemDate(c, &s, card)
		}
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a card",
		Data:    *card,
		Count:   1,
	})
}

// CreateNewCardBulk method
// @Description Create cards
// @Summary create cards
// @Tags Card
// @Produce json
// @Accept json
// @Param card body models.Card true "Cards to create"
// @Success 200
// @Router /v1/cards/deck/{deckID}/bulk [post]
func CreateNewCardBulk(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	deckID, _ := strconv.ParseUint(c.Params("deckID"), 10, 32)

	auth := CheckAuth(c, models.PermMod) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	type Data struct {
		Cards []models.Card `json:"cards"`
	}

	data := new(Data)

	if err := c.BodyParser(&data); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries.CheckAccess(c, auth.User.ID, uint(deckID), models.AccessOwner); !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	for _, card := range data.Cards {
		if len(card.Question) <= 5 || len(card.Answer) <= 5 {
			_ = models.ResponseHTTP{
				Success: false,
				Message: utils.ErrorQALen,
				Data:    nil,
				Count:   0,
			}
		} else {
			card.DeckID = uint(deckID)
			db.Create(card)

			log := queries.CreateLog(models.LogCardCreated, auth.User.Username+" created "+card.Question)
			_ = queries.CreateUserLog(auth.User.ID, log)
			_ = queries.CreateDeckLog(uint(deckID), log)
			_ = queries.CreateCardLog(card.ID, log)

			_ = models.ResponseHTTP{
				Success: true,
				Message: "Success creating card",
				Data:    nil,
				Count:   0,
			}

			var users []models.User

			if users = queries.GetSubUsers(c, uint(deckID)); len(users) > 0 {

				for _, s := range users {
					_ = queries.GenerateMemDate(c, &s, &card)
				}
			}
		}

	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success bulk creation",
		Data:    data.Cards,
		Count:   len(data.Cards),
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
		return queries.AuthError(c, auth)
	}

	response := new(models.CardResponse)
	card := new(models.Card)

	if err := c.BodyParser(&response); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())

	}

	if err := db.Joins("Deck").First(&card, response.CardID).Error; err != nil {
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	res := queries.CheckAccess(c, auth.User.ID, card.Deck.ID, models.AccessStudent)

	if !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
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

	_ = queries.PostMem(c, auth.User, *card, *validation)

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
		return queries.AuthError(c, auth)
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := queries.CheckAccess(c, auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if err := UpdateCard(c, card); !err.Success {
		return queries.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := queries.CreateLog(models.LogCardEdited, auth.User.Username+" edited "+card.Question)
	_ = queries.CreateUserLog(auth.User.ID, log)
	_ = queries.CreateCardLog(card.ID, log)

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
		return queries.AuthError(c, auth)
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	if res := queries.CheckAccess(c, auth.User.ID, card.DeckID, models.AccessOwner); !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	db.Delete(card)

	log := queries.CreateLog(models.LogCardDeleted, auth.User.Username+" deleted "+card.Question)
	_ = queries.CreateUserLog(auth.User.ID, log)
	_ = queries.CreateCardLog(card.ID, log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success delete card by ID",
		Data:    *card,
		Count:   1,
	})

}
