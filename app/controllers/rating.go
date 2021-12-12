package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllRatings method to get all ratings
// @Description Get all ratings. Admin Only
// @Summary get a list of ratings
// @Tags Rating
// @Produce json
// @Success 200 {array} models.Rating
// @Router /v1/ratings [get]
func GetAllRatings(c *fiber.Ctx) error {
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

	var ratings []models.Rating

	if res := db.Joins("User").Find(&ratings); res.Error != nil {

		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all ratings",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all ratings",
		Data:    ratings,
		Count:   len(ratings),
	})
}

// GetAllRatings method to get all ratings
// @Description Get all ratings from a deckID. Admin Only
// @Summary get a list of ratings
// @Tags Rating
// @Produce json
// @Success 200 {array} models.Rating
// @Router /v1/ratings/deck/{deckID} [get]
func GetAllRatingsByDeck(c *fiber.Ctx) error {
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

	deckID := c.Params("deckID")
	var ratings []models.Rating

	if res := db.Joins("User").Where("ratings.deck_id = ? ", deckID).Find(&ratings); res.Error != nil {

		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all ratings",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all ratings",
		Data:    ratings,
		Count:   len(ratings),
	})
}

// GetAverageRatingByDeck method
// @Description Get average rating from a deckID.
// @Summary get an average rating
// @Tags Rating
// @Produce json
// @Success 200 {object} integer
// @Router /v1/ratings/deck/{deckID}/average [get]
func GetAverageRatingByDeck(c *fiber.Ctx) error {
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

	deckID := c.Params("deckID")
	var averageValue float32

	if res := db.Table("ratings").Joins("User").Select("AVG(value)").Where("ratings.deck_id = ? ", deckID).Find(&averageValue); res.Error != nil {

		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get average rating",
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get average rating",
		Data:    averageValue,
		Count:   1,
	})
}

// GetRatingByDeckAndUser method
// @Description Get a rating by user & deck
// @Summary get a rating
// @Tags Rating
// @Produce json
// @Success 200 {object} models.Rating
// @Router /v1/ratings/deck/{deckID}/user/{userID} [get]
func GetRatingByDeckAndUser(c *fiber.Ctx) error {
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

	deckID := c.Params("deckID")
	userID := c.Params("userID")

	rating := new(models.Rating)

	if res := db.Joins("User").Joins("Deck").Where("ratings.deck_id = ? AND ratings.user_id = ?", deckID, userID).First(&rating); res.Error != nil {

		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get a rating",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get a rating",
		Data:    rating,
		Count:   1,
	})
}


// POST


