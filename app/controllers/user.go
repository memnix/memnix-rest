package controllers

import (
	"bytes"
	"fmt"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/app/queries"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET

// GetAllUsers method to get all users
// @Description Get all users.  Shouldn't really be used
// @Summary gets a list of user
// @Tags User
// @Produce json
// @Success 200 {object} models.User
// @Security Admin
// @Deprecated
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
// @Security Admin
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

// SetTodayConfig method to set a config
// @Description Set the today config for a deck
// @Summary sets the today config for a deck
// @Tags User
// @Produce json
// @Accept json
// @Param deckId path int true "Deck ID"
// @Param config body models.DeckConfig true "Deck Config"
// @Success 200
// @Router /v1/users/config/{deckId}/today [post]
func SetTodayConfig(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	// Params
	deckID := c.Params("deckID")
	deckidInt, _ := strconv.ParseUint(deckID, 10, 32)

	deckConfig := new(models.DeckConfig)

	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return queries.AuthError(c, &auth)
	}

	if err := c.BodyParser(&deckConfig); err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on SetTodayConfig: %s from %s", err.Error(), auth.User.Email), models.LogBodyParserError).SetType(models.LogTypeError).AttachIDs(auth.User.ID, uint(deckidInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	access := new(models.Access)
	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", auth.User.ID, deckID).Find(&access).Error; err != nil {
		log := models.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - SetTodayConfig", auth.User.Email, deckidInt), models.LogDeckCardLimit).SetType(models.LogTypeWarning).AttachIDs(auth.User.ID, uint(deckidInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorNotSub)
	}

	if access.Permission == 0 {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorNotSub)
	}

	access.ToggleToday = deckConfig.TodaySetting

	db.Save(access)

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success updated deck config",
		Data:    nil,
		Count:   1,
	})
}

// ResetPassword method to request a password reset
// @Description Request a password reset
// @Summary gets a code to reset a password
// @Tags User
// @Produce json
// @Accept json
// @Param config body string true "Email"
// @Success 200
// @Router /v1/users/resetpassword [post]
func ResetPassword(c *fiber.Ctx) error {
	db := database.DBConn

	var body struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&body); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	email := strings.ToLower(body.Email)
	email = strings.TrimSpace(email)

	if err := db.Where("email = ?", email).First(&models.User{}).Error; err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	_, found := database.Cache.Get(email)
	if found {
		return queries.RequestError(c, http.StatusBadRequest, "A password reset has already been sent to this email address. Please check your email for the link.")
	}

	token := utils.GenerateSecretCode(3)
	database.Cache.Set(email, token, time.Minute*10)

	// Send email
	go func() {
		err := utils.SendEmail(email, "Password Reset", "Your password reset code is: "+token)
		if err != nil {
			log := models.CreateLog(fmt.Sprintf("Error on ResetPassword: %s", err.Error()), models.LogBodyParserError).SetType(models.LogTypeError).AttachIDs(0, 0, 0)
			_ = log.SendLog()
		}
	}()

	log := models.CreateLog(fmt.Sprintf("Password reset request for %s", email), models.LogUserPasswordReset).SetType(models.LogTypeInfo).AttachIDs(0, 0, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success sent password reset email",
		Data:    nil,
		Count:   1,
	})
}

// ResetPasswordConfirm method to confirm a password reset
// @Description Confirm a password reset
// @Summary reset a password
// @Tags User
// @Produce json
// @Accept json
// @Param config body models.PasswordResetConfirm true "Password reset"
// @Success 200
// @Router /v1/users/confirmpassword [post]
func ResetPasswordConfirm(c *fiber.Ctx) error {
	db := database.DBConn

	var body models.PasswordResetConfirm

	if err := c.BodyParser(&body); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	email := strings.ToLower(body.Email)
	email = strings.TrimSpace(email)

	token, found := database.Cache.Get(email)
	if !found {
		return queries.RequestError(c, http.StatusBadRequest, "Your password reset code has expired. Please request a new one.")
	}

	if token != body.Code {
		return queries.RequestError(c, http.StatusBadRequest, "Invalid password reset code.")
	}

	user := new(models.User)
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	// Register checks
	if len(body.Pass) > utils.MaxPasswordLen {
		log := models.CreateLog(fmt.Sprintf("Error on reset password: %s - %s", user.Username, user.Email), models.LogBadRequest).SetType(models.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorRequestFailed)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(body.Pass), 10) // Hash password

	user.Password = password
	db.Save(user)

	database.Cache.Delete(email)

	log := models.CreateLog(fmt.Sprintf("Password reset for %s", email), models.LogUserPasswordChanged).SetType(models.LogTypeInfo).AttachIDs(user.ID, 0, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(models.ResponseHTTP{
		Success: true,
		Message: "Success updated password",
		Data:    nil,
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
