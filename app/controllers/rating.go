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

// GetAllRatings method to get all ratings
// @Description Get all ratings. Admin Only
// @Summary gets a list of ratings
// @Tags Rating
// @Produce json
// @Success 200 {array} models.Rating
// @Router /v1/ratings [get]
func GetAllRatings(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	var ratings []models.Rating

	if res := db.Joins("User").Joins("Deck").Find(&ratings); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all ratings",
		Data:    ratings,
		Count:   len(ratings),
	})
}

// GetAllRatingsByDeck method to get all ratings
// @Description Get all ratings from a deckID. Admin Only
// @Summary gets a list of ratings
// @Tags Rating
// @Produce json
// @Success 200 {array} models.Rating
// @Router /v1/ratings/deck/{deckID} [get]
func GetAllRatingsByDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	deckID := c.Params("deckID")
	var ratings []models.Rating

	if res := db.Joins("User").Joins("Deck").Where("ratings.deck_id = ? ", deckID).Find(&ratings); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, utils.ErrorRequestFailed)
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all ratings",
		Data:    ratings,
		Count:   len(ratings),
	})
}

// GetRatingsByDeck method
// @Description Get ratings from a deckID.
// @Summary get a rating
// @Tags Rating
// @Produce json
// @Success 200 {object} models.Rating
// @Router /v1/ratings/deck/{deckID}/user [get]
func GetRatingsByDeck(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	deckID := c.Params("deckID")
	var ratings []models.Rating

	if res := db.Joins("User").Joins("Deck").Where("ratings.deck_id = ? AND ratings.user_id = ?", deckID, auth.User.ID).First(&ratings); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, res.Error.Error())
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get a rating",
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
		return queries.AuthError(c, auth)
	}

	deckID := c.Params("deckID")
	var averageValue float32

	if res := db.Table("ratings").Select("AVG(value)").Where("ratings.deck_id = ? ", deckID).Find(&averageValue); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, res.Error.Error())
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
// @Summary gets a rating
// @Tags Rating
// @Produce json
// @Success 200 {object} models.Rating
// @Router /v1/ratings/deck/{deckID}/user/{userID} [get]
func GetRatingByDeckAndUser(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	deckID := c.Params("deckID")
	userID := c.Params("userID")

	rating := new(models.Rating)

	if res := db.Joins("User").Joins("Deck").Where("ratings.deck_id = ? AND ratings.user_id = ?", deckID, userID).First(&rating); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, res.Error.Error())
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get a rating",
		Data:    rating,
		Count:   1,
	})
}

// POST

// RateDeck method
// @Description Rate a deck
// @Summary rates a deck
// @Tags Rating
// @Produce json
// @Accept json
// @Param rating body models.Rating true "Rating to create or update"
// @Success 200
// @Router /v1/rating/new [post]
func RateDeck(c *fiber.Ctx) error {

	rating := new(models.Rating)

	// Check auth
	auth := CheckAuth(c, models.PermUser)
	if !auth.Success {
		return queries.AuthError(c, auth)
	}

	if err := c.BodyParser(&rating); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	rating.UserID = auth.User.ID

	if err := queries.GenerateRating(c, rating); !err.Success {
		return queries.RequestError(c, http.StatusInternalServerError, err.Message)
	}

	log := queries.CreateLog(
		models.LogDeckRated, auth.User.Username+" rated "+strconv.FormatUint(
			uint64(rating.DeckID), 10)+" as "+strconv.FormatUint(uint64(rating.Value), 10))
	_ = queries.CreateUserLog(auth.User.ID, log)
	_ = queries.CreateDeckLog(rating.DeckID, log)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success rating the deck",
		Data:    rating,
		Count:   1,
	})
}
