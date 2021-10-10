package handlers

import (
	"memnixrest/database"
	"memnixrest/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllUsers
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DBConn

	var users []models.User

	if res := db.Find(&users); res.Error != nil {

		return c.JSON(ResponseHTTP{
			Success: false,
			Message: "Get All users",
			Data:    nil,
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Get All users",
		Data:    users,
	})

}

// GetUserByID
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get user by ID.",
		Data:    *user,
	})
}

// GetUserByDiscordID
func GetUserByDiscordID(c *fiber.Ctx) error {
	id := c.Params("discordID")
	db := database.DBConn

	user := new(models.User)

	if err := db.Where("users.discord_id = ?", id).First(&user).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get user by ID.",
		Data:    *user,
	})
}

// POST

// CreateNewUser
func CreateNewUser(c *fiber.Ctx) error {
	db := database.DBConn

	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Create(user)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success register an user",
		Data:    *user,
	})
}

// PUT

// UpdateUserByID
func UpdateUserByID(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := UpdateUser(c, user); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Couldn't update the user",
			Data:    nil,
		})
	}

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success update user by Id.",
		Data:    *user,
	})
}

// UpdateUser
func UpdateUser(c *fiber.Ctx, u *models.User) error {
	db := database.DBConn

	if err := c.BodyParser(&u); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	db.Save(u)

	return nil
}
