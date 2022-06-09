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
// @Param deckID path string true "Deck ID"
// @Security Beaver
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

// CreateMcq method
// @Description Create a new mcq
// @Summary creates a mcq
// @Tags Mcq
// @Produce json
// @Accept json
// @Param mcq body models.Mcq true "Mcq to create"
// @Security Beaver
// @Success 200
// @Router /v1/mcqs/new [post]
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

	if res := queries.CheckCardLimit(auth.User.Permissions, mcq.DeckID); !res {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateMcq: This deck has reached his limit", auth.User.Email, mcq.DeckID), models.LogDeckCardLimit).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, "This deck has reached his limit ! You can't add more mcq to it.")
	}

	if mcq.NotValidate() {
		log := models.CreateLog(fmt.Sprintf("Error from %s on CreateMcq: BadRequest", auth.User.Email), models.LogBadRequest).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, "You must provide at least 3 and at most 150 answers for Standalone MCQ")
	}

	db.Create(mcq)

	log := models.CreateLog(fmt.Sprintf("Created MCQ: %d - %s", mcq.ID, mcq.Name), models.LogCardCreated).SetType(models.LogTypeInfo).AttachIDs(auth.User.ID, mcq.DeckID, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a mcq",
		Data:    *mcq,
		Count:   1,
	})
}

// PUT

// UpdateMcqById method
// @Description Edit a mcq
// @Summary edits a mcq
// @Tags Mcq
// @Produce json
// @Success 200
// @Accept json
// @Param mcq body models.Mcq true "MCQ to edit"
// @Param mcqID path string true "MCQ ID"
// @Security Beaver
// @Router /v1/mcqs/{mcqID}/edit [put]
func UpdateMcqById(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	mcq := new(models.Mcq)

	if err := db.First(&mcq, id).Error; err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on UpdateMcqById: %s from %s", err.Error(), auth.User.Email), models.LogQueryGetError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := queries.CheckAccess(auth.User.ID, mcq.DeckID, models.AccessEditor); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - UpdateCardByID: %s", auth.User.Email, mcq.DeckID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if err := UpdateMcq(c, mcq); !err.Success {
		log := models.CreateLog(fmt.Sprintf("Error on UpdateCardByID: %s from %s", err.Message, auth.User.Email), models.LogBadRequest).SetType(models.LogTypeError).AttachIDs(auth.User.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := models.CreateLog(fmt.Sprintf("Edited: %d - %s", mcq.ID, mcq.Name), models.LogCardEdited).SetType(models.LogTypeInfo).AttachIDs(auth.User.ID, mcq.DeckID, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update mcq by ID",
		Data:    *mcq,
		Count:   1,
	})
}

// UpdateMcq function
func UpdateMcq(c *fiber.Ctx, mcq *models.Mcq) *models.ResponseHTTP {
	db := database.DBConn

	deckId := mcq.DeckID

	res := new(models.ResponseHTTP)

	if err := c.BodyParser(&mcq); err != nil {
		res.GenerateError(err.Error())
		return res
	}

	if deckId != mcq.DeckID {
		res.GenerateError(utils.ErrorBreak)
		return res
	}

	if mcq.Type == models.McqStandalone && (len(mcq.Answers) < utils.MinMcqAnswersLen || len(mcq.Answers) > utils.MaxMcqAnswersLen) || len(mcq.Name) > utils.MaxMcqName || mcq.Name == "" {
		res.GenerateError(utils.ErrorRequestFailed)
		return res
	}

	if mcq.Type == models.McqLinked {
		mcq.UpdateLinkedAnswers()
	}

	db.Save(mcq)

	res.GenerateSuccess("Success update mcq", nil, 0)
	return res
}

// DeleteMcqById method
// @Description Delete a mcq
// @Summary deletes a mcq
// @Tags Mcq
// @Produce json
// @Success 200
// @Param mcqID path string true "MCQ ID"
// @Security Beaver
// @Router /v1/mcqs/{mcqID} [delete]
func DeleteMcqById(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn
	id := c.Params("id")

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	mcq := new(models.Mcq)

	if err := db.First(&mcq, id).Error; err != nil {
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	if res := queries.CheckAccess(auth.User.ID, mcq.DeckID, models.AccessOwner); !res.Success {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d and mcq %d - DeleteMcqById: %s", auth.User.Email, mcq.DeckID, mcq.ID, res.Message), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	db.Delete(mcq)

	log := models.CreateLog(fmt.Sprintf("Deleted: %d - %s", mcq.ID, mcq.Name), models.LogCardDeleted).SetType(models.LogTypeInfo).AttachIDs(auth.User.ID, mcq.DeckID, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success delete mcq by ID",
		Data:    *mcq,
		Count:   1,
	})

}
