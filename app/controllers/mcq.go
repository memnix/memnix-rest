package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/app/queries"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"net/http"
	"strconv"
)

// GetMcqsByDeck method
// @Description Get mcqs linked to the deck
// @Summary gets a list of mcqs
// @Tags Mcq
// @Produce json
// @Success 200 {array} models.Mcq
// @Router /v1/mcqs/{deckID} [get]
func GetMcqsByDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	// Params
	deckID := c.Params("deckID")

	var mcqs []models.Mcq

	if err := db.Joins("Deck").Where("mcqs.deck_id = ?", deckID).Find(&mcqs).Error; err != nil {
		deckidInt, _ := strconv.ParseUint(deckID, 10, 32)
		log := models.CreateLog(fmt.Sprintf("Error from %s on GetMcqsByDeck: %s", auth.User.Email, err.Error()), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, uint(deckidInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get mcqs by deck.",
		Data:    mcqs,
		Count:   len(mcqs),
	})
}

func CreateMcq(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	mcq := new(models.Mcq)

	if err := c.BodyParser(&mcq); err != nil {
		log := models.CreateLog(fmt.Sprintf("Error from %s on CreateMcq: %s", auth.User.Email, err.Error()), models.LogBodyParserError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries.CheckAccess(auth.User.ID, mcq.DeckID, models.AccessEditor); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateMcq: %s", auth.User.Email, mcq.DeckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if mcq.Type == models.McqStandalone && len(mcq.Answers) < 3 {
		log := models.CreateLog(fmt.Sprintf("Error from %s on CreateMcq: BadRequest", auth.User.Email), models.LogBadRequest).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, "You must provide at least 3 answers for Standalone MCQ")
	}

	db.Create(mcq)

	//TODO: Log MCQ

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a mcq",
		Data:    *mcq,
		Count:   1,
	})
}
