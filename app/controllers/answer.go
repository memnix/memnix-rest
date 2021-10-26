package handlers

import (
	"memnixrest/database"
	"memnixrest/app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllAnswers
func GetAllAnswers(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var answers []models.Answer

	if res := db.Joins("Card").Find(&answers); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get All answers",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All answers",
		Data:    answers,
		Count:   len(answers),
	})

}

// GetAnswerByID
func GetAnswerByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	answer := new(models.Answer)

	if err := db.Joins("Card").First(&answer, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get answer by ID.",
		Data:    *answer,
		Count:   1,
	})
}

// GetAnswersByCardID
func GetAnswersByCardID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	cardID := c.Params("cardID")

	var answers []models.Answer

	if res := db.Joins("Card").Where("answers.card_id = ?", cardID).Find(&answers); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get answers by cardID",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get answers by cardID",
		Data:    answers,
		Count:   len(answers),
	})
}

// POST

// CreateNewAnswer
func CreateNewAnswer(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	answer := new(models.Answer)

	if err := c.BodyParser(&answer); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Preload("Card").Create(answer)

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register an answer",
		Data:    *answer,
		Count:   1,
	})

}
