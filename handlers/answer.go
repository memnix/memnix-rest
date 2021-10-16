package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllAnswers
func GetAllAnswers(c *fiber.Ctx) error {
	db := database.DBConn

	var answers []models.Answer

	if res := db.Joins("Card").Find(&answers); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get All answers",
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
	id := c.Params("id")
	db := database.DBConn

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
	db := database.DBConn

	cardID := c.Params("cardID")

	var answers []models.Access

	if res := db.Joins("Card").Where("answers.card_id = ?", cardID).Find(&answers); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get all answers",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all answers",
		Data:    answers,
		Count:   len(answers),
	})
}

// POST

// CreateNewAnswer
func CreateNewAnswer(c *fiber.Ctx) error {
	db := database.DBConn

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
		Message: "Success register a answer",
		Data:    *answer,
		Count:   1,
	})

}
