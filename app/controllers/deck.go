package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"memnixrest/pkg/queries"
	"memnixrest/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllDecks method to get all decks
// @Description Get every deck. Shouldn't really be used, consider using /v1/decks/public instead !
// @Summary gets all decks
// @Tags Deck
// @Produce json
// @Success 200 {object} models.Deck
// @Router /v1/decks [get]
func GetAllDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	var decks []models.Deck

	if res := db.Find(&decks); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all decks",
		Data:    decks,
		Count:   len(decks),
	})

}

// GetDeckByID method to get a deck
// @Description Get a deck by tech ID
// @Summary get a deck
// @Tags Deck
// @Produce json
// @Param id path int true "Deck ID"
// @Success 200 {model} models.Deck
// @Router /v1/decks/{deckID} [get]
func GetDeckByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("deckID")

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	deck := new(models.Deck)

	if err := db.First(&deck, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get deck by ID.",
		Data:    *deck,
		Count:   1,
	})
}

// GetAllSubDecks method to get a deck
// @Description Get decks a user is sub to
// @Summary get a list of deck
// @Tags Deck
// @Produce json
// @Success 200 {array} models.ResponseDeck
// @Router /v1/decks/sub [get]
func GetAllSubDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	var responseDeck []models.ResponseDeck

	var accesses []models.Access // Accesses array

	if err := db.Joins("Deck").Joins("User").Where("accesses.user_id = ? AND accesses.permission >= ?", auth.User.ID, models.AccessStudent).Find(&accesses).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	for _, s := range accesses {
		responseDeck = append(responseDeck, queries.FillResponseDeck(c, &s))
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all sub decks",
		Data:    responseDeck,
		Count:   len(responseDeck),
	})
}

// GetAllSubUsers method to get a list of users
// @Description Get all the sub users to a deck
// @Summary gets a list of users
// @Tags Deck
// @Produce json
// @Success 200 {array} models.User
// @Router /v1/decks/{deckID}/users [get]
func GetAllSubUsers(c *fiber.Ctx) error {

	// Params
	deckID := c.Params("deckID")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	var users []models.User
	id, _ := strconv.ParseUint(deckID, 10, 32)

	if users = queries.GetSubUsers(c, uint(id)); len(users) == 0 || users == nil {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all sub users",
		Data:    users,
		Count:   len(users),
	})
}

// GetAllAvailableDecks method to get a list of deck
// @Description Get all public deck that you are not sub to
// @Summary get a list of deck
// @Tags Deck
// @Produce json
// @Success 200 {array} models.ResponseDeck
// @Router /v1/decks/available [get]
func GetAllAvailableDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	var responseDeck []models.ResponseDeck

	var accesses []models.Access // Accesses array

	if err := db.Joins("join decks on accesses.deck_id = decks.id").Joins("Deck").Joins("User").Where("decks.status = ? AND accesses.user_id = ? AND((accesses.deck_id IS NULL) OR (accesses.permission < ?))", models.DeckPublic, auth.User.ID, models.AccessStudent).Find(&accesses).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	for _, s := range accesses {
		responseDeck = append(responseDeck, queries.FillResponseDeck(c, &s))
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all available decks",
		Data:    responseDeck,
		Count:   len(responseDeck),
	})
}

// GetAllPublicDecks method to get a list of deck
// @Description Get all public deck
// @Summary gets a list of deck
// @Tags Deck
// @Produce json
// @Success 200 {array} models.Deck
// @Router /v1/decks/public [get]
func GetAllPublicDecks(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	var decks []models.Deck

	if err := db.Where("decks.status = ?", models.DeckPublic).Find(&decks).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all public decks",
		Data:    decks,
		Count:   len(decks),
	})
}

// POST

// CreateNewDeck method
// @Description Create a new deck
// @Summary creates a deck
// @Tags Deck
// @Produce json
// @Success 200
// @Accept json
// @Param deck body models.Deck true "Deck to create"
// @Router /v1/decks/new [post]
func CreateNewDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	deck := new(models.Deck)

	if err := c.BodyParser(&deck); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())

	}

	if len(deck.DeckName) <= 5 {
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorDeckName)
	}

	deck.Status = models.DeckPrivate
	db.Create(deck)

	if err := queries.GenerateCreatorAccess(c, &auth.User, deck); !err.Success {
		return queries.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := queries.CreateLog(models.LogDeckCreated, auth.User.Username+" created "+deck.DeckName)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateDeckLog(deck.ID, *log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a deck",
		Data:    *deck,
		Count:   1,
	})
}

// UnSubToDeck method
// @Description Unsubscribe to a deck
// @Summary unsub deck
// @Tags Deck
// @Produce json
// @Success 200
// @Accept json
// @Router /v1/decks/{deckID}/unsubscribe [post]
func UnSubToDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	// Params
	deckID := c.Params("deckID")

	access := new(models.Access)
	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", auth.User.ID, deckID).Find(&access).Error; err != nil {
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorNotSub)
	}

	if access.Permission == 0 {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorNotSub)
	}

	access.Permission = 0
	db.Preload("User").Preload("Deck").Save(access)

	_ = queries.DeleteRating(c, &auth.User, &access.Deck)

	log := queries.CreateLog(models.LogUnsubscribe, auth.User.Username+" unsubscribed to "+access.Deck.DeckName)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateDeckLog(access.Deck.ID, *log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success unsub to the deck",
		Data:    nil,
		Count:   0,
	})
}

// SubToDeck method
// @Description Subscribe to a deck
// @Summary sub deck
// @Tags Deck
// @Produce json
// @Success 200
// @Accept json
// @Router /v1/decks/{deckID}/subscribe [post]
func SubToDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	deckID := c.Params("deckID")
	deck := new(models.Deck)

	// Check auth
	auth := CheckAuth(c, models.PermUser)
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	if err := db.First(&deck, deckID).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())

	}

	if err := queries.GenerateAccess(c, &auth.User, deck); !err.Success {
		return queries.RequestError(c, http.StatusInternalServerError, err.Message)
	}

	if err := queries.PopulateMemDate(c, &auth.User, deck); !err.Success {
		return queries.RequestError(c, http.StatusInternalServerError, err.Message)
	}

	log := queries.CreateLog(models.LogSubscribe, auth.User.Username+" subscribed to "+deck.DeckName)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateDeckLog(deck.ID, *log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success subscribing to deck",
		Data:    nil,
		Count:   0,
	})
}

// PUT

// UpdateDeckByID method
// @Description Edit a deck
// @Summary edits a deck
// @Tags Deck
// @Produce json
// @Success 200
// @Accept json
// @Param deck body models.Deck true "Deck to edit"
// @Router /v1/decks/{deckID}/edit [put]
func UpdateDeckByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("deckID")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	deck := new(models.Deck)

	if err := db.First(&deck, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())

	}

	if res := queries.CheckAccess(c, auth.User.ID, deck.ID, models.AccessOwner); !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)

	}

	if err := UpdateDeck(c, deck); !err.Success {
		return queries.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := queries.CreateLog(models.LogDeckEdited, auth.User.Username+" edited "+deck.DeckName)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateDeckLog(deck.ID, *log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update deck by ID",
		Data:    *deck,
		Count:   1,
	})
}

// UpdateDeck function
func UpdateDeck(c *fiber.Ctx, d *models.Deck) *models.ResponseHTTP {
	db := database.DBConn

	deckStatus := d.Status

	res := new(models.ResponseHTTP)

	if err := c.BodyParser(&d); err != nil {
		res.GenerateError(err.Error())
		return res

	}

	if d.Status != deckStatus {
		res.GenerateError(utils.ErrorBreak)
		return res
	}

	if len(d.DeckName) <= 5 {
		res.GenerateError(utils.ErrorQALen)
		return res
	}

	db.Save(d)

	res.GenerateSuccess("Success update deck", nil, 0)
	return res
}

// DeleteDeckById method
// @Description Delete a deck
// @Summary delete a deck
// @Tags Deck
// @Produce json
// @Success 200
// @Router /v1/decks/{deckID} [delete]
func DeleteDeckById(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	id := c.Params("deckID")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	deck := new(models.Deck)

	if err := db.First(&deck, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())

	}

	if res := queries.CheckAccess(c, auth.User.ID, deck.ID, models.AccessOwner); !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	db.Delete(deck)

	log := queries.CreateLog(models.LogDeckDeleted, auth.User.Username+" deleted "+deck.DeckName)
	_ = queries.CreateUserLog(auth.User.ID, *log)
	_ = queries.CreateDeckLog(deck.ID, *log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success delete deck by ID",
		Data:    *deck,
		Count:   1,
	})

}
