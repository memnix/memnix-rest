package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/app/queries"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"net/http"
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
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries.CheckAccess(auth.User.ID, mcq.DeckID, models.AccessEditor); !res.Success {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if mcq.Type == models.McqStandalone && len(mcq.Answers) < 3 {
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
