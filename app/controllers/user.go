package controllers

import (
	"bytes"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/app/queries"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllUsers method to get all users
// @Description Get all users.  Shouldn't really be used
// @Summary gets a list of user
// @Tags User
// @Produce json
// @Success 200 {object} models.User
// @Router /v1/users [get]
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	var users []models.User

	if res := db.Find(&users); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, res.Error.Error())
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Get all users",
		Data:    users,
		Count:   len(users),
	})
}

// GetUserByID method to get a user
// @Description Get a user by ID.
// @Summary gets a user
// @Tags User
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Router /v1/users/id/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	// Params
	id := c.Params("id")

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success get user by ID.",
		Data:    *user,
		Count:   1,
	})
}

// PUT

// UpdateUserByID function
func UpdateUserByID(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	id := c.Params("id")

	auth := CheckAuth(c, models.PermAdmin) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := UpdateUser(c, user); !res.Success {
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success update user by ID",
		Data:    *user,
		Count:   1,
	})
}

// UpdateUser function
func UpdateUser(c *fiber.Ctx, u *models.User) *models.ResponseHTTP {
	db := database.DBConn

	email, password, permissions := u.Email, u.Password, u.Permissions

	res := new(models.ResponseHTTP)

	if err := c.BodyParser(&u); err != nil {
		res.GenerateError(err.Error())
		return res
	}

	if u.Email != email || !bytes.Equal(u.Password, password) || u.Permissions != permissions {
		res.GenerateError(utils.ErrorBreak)
		return res
	}

	db.Save(u)

	res.GenerateSuccess("Success update user", nil, 0)
	return res
}
