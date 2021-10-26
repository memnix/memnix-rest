package handlers

import (
	"memnixrest/app/models"
	"memnixrest/database"
	"memnixrest/pkg/core"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllMem
func GetAllMem(c *fiber.Ctx) error {
	db := database.DBConn

	var mems []models.Mem

	if res := db.Find(&mems); res.Error != nil {

		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Get All mems",
			Data:    nil,
			Count:   0,
		})
	}
	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get All mems",
		Data:    mems,
		Count:   len(mems),
	})

}

// GetMemByID
func GetMemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	mem := new(models.Mem)

	if err := db.Joins("User").Joins("Card").First(&mem, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get mem by ID.",
		Data:    *mem,
		Count:   1,
	})
}

// GetMemByCardAndUser
func GetMemByCardAndUser(c *fiber.Ctx) error {
	userID := c.Params("userID")
	cardID := c.Params("cardID")

	db := database.DBConn

	mem := new(models.Mem)

	if err := db.Joins("User").Joins("Card").Where("mems.user_id = ? AND mems.card_id = ?", userID, cardID).First(&mem).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get mem by UserID & CardID.",
		Data:    *mem,
		Count:   1,
	})
}

// GetNextMem
func GetNextMem(c *fiber.Ctx) error {
	userIDTemp := c.Params("userID")
	deckIDTemp := c.Params("deckID")

	userID, _ := strconv.Atoi(userIDTemp)
	deckID, _ := strconv.Atoi(deckIDTemp)

	res := core.FetchNextCard(c, uint(userID), uint(deckID))

	if !res.Success {
		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "Next card not found",
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get card by ID.",
		Data:    res.Data,
		Count:   1,
	})

}

// GetTodayNextMem
func GetTodayNextMem(c *fiber.Ctx) error {
	userIDTemp := c.Params("userID")
	deckIDTemp := c.Params("deckID")

	userID, _ := strconv.Atoi(userIDTemp)
	deckID, _ := strconv.Atoi(deckIDTemp)

	res := core.FetchNextTodayCard(c, uint(userID), uint(deckID))
	if !res.Success {
		return c.JSON(models.ResponseHTTP{
			Success: false,
			Message: "No more card for today!",
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get card Today's card.",
		Data:    res.Data,
		Count:   1,
	})
}

// POST

// CreateNewMem
func CreateNewMem(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	mem := new(models.Mem)

	if err := c.BodyParser(&mem); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Preload("User").Preload("Card").Create(mem)

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register a new mem",
		Data:    *mem,
		Count:   1,
	})
}

// PUT

// UpdateMemByID
func UpdateMemByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	mem := new(models.Mem)

	if err := db.First(&mem, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if err := UpdateMem(c, mem); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Couldn't update the mem",
			Data:    nil,
			Count:   0,
		})
	}

	return c.JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update mem by Id.",
		Data:    *mem,
		Count:   1,
	})
}

// UpdateMem
func UpdateMem(c *fiber.Ctx, m *models.Mem) error {
	db := database.DBConn // DB Conn

	if err := c.BodyParser(&m); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Save(m)

	return nil
}
