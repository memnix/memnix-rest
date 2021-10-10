package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllMem
func GetAllMem(c *fiber.Ctx) error {
	db := database.DBConn

	var mems []models.Mem

	if res := db.Find(&mems); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All mems",
			Data:    nil,
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All mems",
		Data:    mems,
	})

}

// GetMemByID
func GetMemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	mem := new(models.Mem)

	if err := db.Joins("User").Joins("Card").First(&mem, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get mem by ID.",
		Data:    *mem,
	})
}

// GetMemByCardAndUser
func GetMemByCardAndUser(c *fiber.Ctx) error {
	userID := c.Params("userID")
	cardID := c.Params("cardID")

	db := database.DBConn

	mem := new(models.Mem)

	if err := db.Joins("User").Joins("Card").Where("mems.user_id = ? AND mems.card_id = ?", userID, cardID).First(&mem).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get mem by UserID & CardID.",
		Data:    *mem,
	})
}

// POST

// SubToDeck
func SubToDeck(c *fiber.Ctx) error {
	id := c.Params("deckID")
	db := database.DBConn

	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", id).Find(&cards).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	for x := 0; x < len(cards); x++ {
		mem := new(models.Mem)

		if err := c.BodyParser(&mem); err != nil {
			return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
		}

		mem.CardID = cards[x].ID

		db.Preload("User").Preload("Card").Create(mem)

	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success subscribing to deck",
		Data:    nil,
	})
}

// CreateNewMem
func CreateNewMem(c *fiber.Ctx) error {
	db := database.DBConn

	mem := new(models.Mem)

	if err := c.BodyParser(&mem); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Preload("User").Preload("Card").Create(mem)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success subscribing to deck",
		Data:    *mem,
	})
}

// PUT

// UpdateMemByID
func UpdateMemByID(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	mem := new(models.Mem)

	if err := db.First(&mem, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := UpdateMem(c, mem); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Couldn't update the mem",
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success update mem by Id.",
		Data:    *mem,
	})
}

// UpdateMem
func UpdateMem(c *fiber.Ctx, m *models.Mem) error {
	db := database.DBConn

	if err := c.BodyParser(&m); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Save(m)

	return nil
}
