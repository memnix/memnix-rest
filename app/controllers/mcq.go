package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/pkg/logger"
	"github.com/memnix/memnixrest/pkg/queries"
	"github.com/memnix/memnixrest/utils"
	"github.com/memnix/memnixrest/viewmodels"
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
	db := infrastructures.GetDBConn() // DB Conn

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	// Params
	deckID := c.Params("deckID")

	var mcqs []models.Mcq

	if err := db.Joins("Deck").Where("mcqs.deck_id = ?", deckID).Find(&mcqs).Error; err != nil {
		deckidInt, _ := strconv.ParseUint(deckID, 10, 32)
		log := logger.CreateLog(fmt.Sprintf("Error from %s on GetMcqsByDeck: %s", user.Email, err.Error()), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, uint(deckidInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
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
	db := infrastructures.GetDBConn() // DB Conn

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	mcq := new(models.Mcq)

	if err := c.BodyParser(&mcq); err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error from %s on CreateMcq: %s", user.Email, err.Error()), logger.LogBodyParserError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	if res := queries.CheckAccess(user.ID, mcq.DeckID, models.AccessEditor); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateMcq: %s", user.Email, mcq.DeckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if res := queries.CheckCardLimit(user.Permissions, mcq.DeckID); !res {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - CreateMcq: This deck has reached his limit", user.Email, mcq.DeckID), logger.LogDeckCardLimit).SetType(logger.LogTypeWarning).AttachIDs(user.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, "This deck has reached his limit ! You can't add more mcq to it.")
	}

	if mcq.NotValidate() {
		log := logger.CreateLog(fmt.Sprintf("Error from %s on CreateMcq: BadRequest", user.Email), logger.LogBadRequest).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, "You must provide at least 3 and at most 150 answers for Standalone MCQ")
	}

	db.Create(mcq)

	log := logger.CreateLog(fmt.Sprintf("Created MCQ: %d - %s", mcq.ID, mcq.Name), logger.LogCardCreated).SetType(logger.LogTypeInfo).AttachIDs(user.ID, mcq.DeckID, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
		Success: true,
		Message: "Success register a mcq",
		Data:    *mcq,
		Count:   1,
	})
}

// PUT

// UpdateMcqByID method
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
func UpdateMcqByID(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn() // DB Conn

	// Params
	id := c.Params("id")

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	mcq := new(models.Mcq)

	if err := db.First(&mcq, id).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on UpdateMcqById: %s from %s", err.Error(), user.Email), logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := queries.CheckAccess(user.ID, mcq.DeckID, models.AccessEditor); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - UpdateCardByID: %s", user.Email, mcq.DeckID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	if err := UpdateMcq(c, mcq); !err.Success {
		log := logger.CreateLog(fmt.Sprintf("Error on UpdateCardByID: %s from %s", err.Message, user.Email), logger.LogBadRequest).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Message)
	}

	log := logger.CreateLog(fmt.Sprintf("Edited: %d - %s", mcq.ID, mcq.Name), logger.LogCardEdited).SetType(logger.LogTypeInfo).AttachIDs(user.ID, mcq.DeckID, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
		Success: true,
		Message: "Success update mcq by ID",
		Data:    *mcq,
		Count:   1,
	})
}

// UpdateMcq function
func UpdateMcq(c *fiber.Ctx, mcq *models.Mcq) *viewmodels.ResponseHTTP {
	db := infrastructures.GetDBConn()

	deckID := mcq.DeckID

	res := new(viewmodels.ResponseHTTP)

	if err := c.BodyParser(&mcq); err != nil {
		res.GenerateError(err.Error())
		return res
	}

	if deckID != mcq.DeckID {
		res.GenerateError(utils.ErrorBreak)
		return res
	}

	if mcq.Type == models.McqStandalone && (len(mcq.Answers) < utils.MinMcqAnswersLen || len(mcq.Answers) > utils.MaxMcqAnswersLen) || len(mcq.Name) > utils.MaxMcqName || mcq.Name == "" {
		res.GenerateError(utils.ErrorRequestFailed)
		return res
	}

	if mcq.Type == models.McqLinked {
		viewmodels.UpdateLinkedAnswers(mcq)
	}

	db.Save(mcq)

	res.GenerateSuccess("Success update mcq", nil, 0)
	return res
}

// DeleteMcqByID method
// @Description Delete a mcq
// @Summary deletes a mcq
// @Tags Mcq
// @Produce json
// @Success 200
// @Param mcqID path string true "MCQ ID"
// @Security Beaver
// @Router /v1/mcqs/{mcqID} [delete]
func DeleteMcqByID(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn() // DB Conn
	id := c.Params("id")

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	mcq := new(models.Mcq)

	if err := db.First(&mcq, id).Error; err != nil {
		return queries.RequestError(c, http.StatusServiceUnavailable, err.Error())
	}

	if res := queries.CheckAccess(user.ID, mcq.DeckID, models.AccessOwner); !res.Success {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d and mcq %d - DeleteMcqById: %s", user.Email, mcq.DeckID, mcq.ID, res.Message), logger.LogPermissionForbidden).SetType(logger.LogTypeWarning).AttachIDs(user.ID, mcq.DeckID, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorForbidden)
	}

	db.Delete(mcq)

	log := logger.CreateLog(fmt.Sprintf("Deleted: %d - %s", mcq.ID, mcq.Name), logger.LogCardDeleted).SetType(logger.LogTypeInfo).AttachIDs(user.ID, mcq.DeckID, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
		Success: true,
		Message: "Success delete mcq by ID",
		Data:    *mcq,
		Count:   1,
	})
}
