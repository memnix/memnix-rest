package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllUsers method to get all users
// @Description Get all users.  Shouldn't really be used
// @Summary get a list of user
// @Tags User
// @Produce json
// @Success 200 {object} models.User
// @Router /v1/users [get]
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var users []models.User

	if res := db.Find(&users); res.Error != nil {

		return c.Status(http.StatusInternalServerError).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Failed to get all users",
			Data:    nil,
			Count:   0,
		})
	}
	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all users",
		Data:    users,
		Count:   len(users),
	})

}

// GetUserByID method to get an user
// @Description Get an user by ID.
// @Summary get an user
// @Tags User
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Router /v1/users/id/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get user by ID.",
		Data:    *user,
		Count:   1,
	})
}

// GetUserByDiscordID
func GetUserByDiscordID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("discordID")

	user := new(models.User)

	if err := db.Where("users.discord_id = ?", id).First(&user).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get user by discordID.",
		Data:    *user,
		Count:   1,
	})
}

// POST

// CreateNewUser
func CreateNewUser(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if err := db.Where("users.discord_id = ?", user.DiscordID).First(&user).Error; err != nil {
		db.Create(user)
	} else {
		db.Save(user)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success register an user",
		Data:    *user,
		Count:   1,
	})
}

// PUT

// UpdateUserByID
func UpdateUserByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	if err := UpdateUser(c, user); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(models.ResponseHTTP{
			Success: false,
			Message: "Couldn't update the user",
			Data:    nil,
			Count:   0,
		})
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update user by ID",
		Data:    *user,
		Count:   1,
	})
}

// UpdateUser
func UpdateUser(c *fiber.Ctx, u *models.User) error {
	db := database.DBConn

	if err := c.BodyParser(&u); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		})
	}

	db.Save(u)

	return nil
}
