package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"memnixrest/pkg/queries"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetTodayCard method
// @Description Get next today card
// @Summary get a card
// @Tags Card
// @Produce json
// @Success 200 {object} models.Card
// @Router /v1/cards/today [get]
func GetTodayCard(c *fiber.Ctx) error {
	res := *new(models.ResponseHTTP)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	if res = queries.FetchNextTodayCard(c, &auth.User); !res.Success {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: res.Message,
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get todays card",
		Data:    res.Data,
		Count:   1,
	})
}

// GetNextCard method
// @Description Get next card
// @Summary get a card
// @Tags Card
// @Produce json
// @Success 200 {object} models.Card
// @Router /v1/cards/next [get]
func GetNextCard(c *fiber.Ctx) error {
	//db := database.DBConn // DB Conn

	res := *new(models.ResponseHTTP)
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	if res = queries.FetchNextCard(c, &auth.User); !res.Success {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: res.Message,
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get next card",
		Data:    res.Data,
		Count:   1,
	})
}

// GetNextCard method
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
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	if res = queries.FetchNextCardByDeck(c, &auth.User, deckID); !res.Success {
		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: res.Message,
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get next card by deck",
		Data:    res.Data,
		Count:   1,
	})
}

// GetAllCards method
// @Description Get every cards. Shouldn't really be used
// @Summary get all cards
// @Tags Card
// @Produce json
// @Success 200 {array} models.Card
// @Router /v1/cards/ [get]
func GetAllCards(c *fiber.Ctx) error {
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

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

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

// CreateNewCard method
// @Description Create a new card
// @Summary create a card
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
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	if err := c.BodyParser(&card); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if res := queries.CheckAccess(c, auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permission to add a card to this deck !",
			Data:    nil,
			Count:   0,
		})
	}

	if len(card.Question) < 1 || len(card.Answer) < 1 {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: "You must provide a question and an answer.",
			Data:    nil,
			Count:   0,
		})
	}

	db.Create(card)

	log := queries.CreateLog(models.LogCardCreated, auth.User.Username+" created "+card.Question)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateDeckLog(card.DeckID, *log)
	_ = queries.CreateCardLog(card.ID, *log)

	var users []models.User

	if users = queries.GetSubUsers(c, card.DeckID); len(users) > 0 {
		ch := make(chan models.ResponseHTTP)

		for _, s := range users {
			go func(c *fiber.Ctx, user models.User, card *models.Card) {
				res := queries.GenerateMemDate(c, &user, card)
				ch <- res
			}(c, s, card)
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

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	type Data struct {
		Cards []models.Card `json:"cards"`
	}

	data := new(Data)

	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if res := queries.CheckAccess(c, auth.User.ID, uint(deckID), models.AccessEditor); !res.Success {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permission to add a card to this deck !",
			Data:    nil,
			Count:   0,
		})
	}

	ch := make(chan models.ResponseHTTP)

	for _, card := range data.Cards {
		go func(c *fiber.Ctx, card models.Card, deckID uint) {

			print("Card: " + card.Question)

			var res models.ResponseHTTP

			if len(card.Question) < 1 || len(card.Answer) < 1 {
				res = models.ResponseHTTP{
					Success: false,
					Message: "You must provide a question and an answer.",
					Data:    nil,
					Count:   0,
				}
			} else {
				ca := &card
				ca.DeckID = deckID
				db.Create(ca)

				log := queries.CreateLog(models.LogCardCreated, auth.User.Username+" created "+ca.Question)
				_ = queries.CreateUserLog(auth.User.ID, *log)
				_ = queries.CreateDeckLog(deckID, *log)
				_ = queries.CreateCardLog(ca.ID, *log)

				res = models.ResponseHTTP{
					Success: true,
					Message: "Success creating card",
					Data:    nil,
					Count:   0,
				}

				var users []models.User

				if users = queries.GetSubUsers(c, deckID); len(users) > 0 {
					ch2 := make(chan models.ResponseHTTP)

					for _, s := range users {
						go func(c *fiber.Ctx, user models.User, card models.Card) {
							res2 := queries.GenerateMemDate(c, &user, &card)
							ch2 <- res2
						}(c, s, card)
					}
				}
			}
			ch <- res
		}(c, card, uint(deckID))
	}
	//TODO: handle errors in chan

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success bulk creation",
		Data:    data.Cards,
		Count:   len(data.Cards),
	})
}

// PostResponse method
// @Description Post a response and check it
// @Summary post a response
// @Tags Card
// @Produce json
// @Success 200
// @Accept json
// @Router /v1/cards/response [post]
func PostResponse(c *fiber.Ctx) error {
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

	response := new(models.CardResponse)
	card := new(models.Card)

	if err := c.BodyParser(&response); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if err := db.Joins("Deck").First(&card, response.CardID).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	res := queries.CheckAccess(c, auth.User.ID, card.Deck.ID, models.AccessStudent)

	if !res.Success {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permission to answer this deck!",
			Data:    nil,
			Count:   0,
		})
	}

	validation := new(models.CardResponseValidation)

	if strings.EqualFold(
		strings.Replace(response.Response, " ", "", -1), strings.Replace(card.Answer, " ", "", -1)) {
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
// @Summary edit a card
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
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if res := queries.CheckAccess(c, auth.User.ID, card.DeckID, models.AccessEditor); !res.Success {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permission to edit this card!",
			Data:    nil,
			Count:   0,
		})
	}

	if err := UpdateCard(c, card); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Couldn't update the card",
			Data:    nil,
			Count:   0,
		})
	}

	log := queries.CreateLog(models.LogCardEdited, auth.User.Username+" edited "+card.Question)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateCardLog(card.ID, *log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update card by ID",
		Data:    *card,
		Count:   1,
	})
}

// UpdateDeck
func UpdateCard(c *fiber.Ctx, card *models.Card) error {
	db := database.DBConn

	if err := c.BodyParser(&card); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if len(card.Question) < 1 || len(card.Answer) < 1 {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: "You must provide a question and an answer",
			Data:    nil,
			Count:   0,
		})
	}

	db.Save(card)

	return nil
}

// DeleteCardById method
// @Description Delete a card
// @Summary delete a card
// @Tags Card
// @Produce json
// @Success 200
// @Router /v1/cards/{cardID} [delete]
func DeleteCardById(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	id := c.Params("id")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	card := new(models.Card)

	if err := db.First(&card, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if res := queries.CheckAccess(c, auth.User.ID, card.DeckID, models.AccessOwner); !res.Success {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permission to delete this card!",
			Data:    nil,
			Count:   0,
		})
	}

	db.Delete(card)

	log := queries.CreateLog(models.LogCardDeleted, auth.User.Username+" deleted "+card.Question)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateCardLog(card.ID, *log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success delete card by ID",
		Data:    *card,
		Count:   1,
	})

}
