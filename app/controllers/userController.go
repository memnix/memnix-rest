package controllers

import (
	"bytes"
	"fmt"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/interfaces"
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/pkg/logger"
	"github.com/memnix/memnixrest/pkg/queries"
	"github.com/memnix/memnixrest/utils"
	"github.com/memnix/memnixrest/viewmodels"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	interfaces.IUserService
}

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
func (u *UserController) GetAllUsers(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn() // DB Conn

	var users []models.User

	if res := db.Find(&users); res.Error != nil {
		return queries.RequestError(c, http.StatusInternalServerError, res.Error.Error())
	}

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
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
func (u *UserController) GetUserByID(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn() // DB Conn

	// Params
	id := c.Params("id")

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
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
// @Router /v1/users/settings/{deckId}/today [post]
func (u *UserController) SetTodayConfig(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn() // DB Conn

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	// Params
	deckID := c.Params("deckID")
	deckidInt, _ := strconv.ParseUint(deckID, 10, 32)

	deckConfig := new(viewmodels.DeckConfig)

	if err := c.BodyParser(&deckConfig); err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on SetTodayConfig: %s from %s", err.Error(), user.Email), logger.LogBodyParserError).SetType(logger.LogTypeError).AttachIDs(user.ID, uint(deckidInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	access := new(models.Access)
	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", user.ID, deckID).Find(&access).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Forbidden from %s on deck %d - SetTodayConfig", user.Email, deckidInt), logger.LogDeckCardLimit).SetType(logger.LogTypeWarning).AttachIDs(user.ID, uint(deckidInt), 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusBadRequest, utils.ErrorNotSub)
	}

	if access.Permission == 0 {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorNotSub)
	}

	access.ToggleToday = deckConfig.TodaySetting
	db.Save(access)
	if deckConfig.TodaySetting {
		_, _ = queries.FetchTodayMemDateByDeck(user.ID, uint(deckidInt), true)
	} else {
		_ = queries.ClearCacheByUserID(user.ID)
	}
	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
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
func (u *UserController) ResetPassword(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn()

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

	_, found := infrastructures.Cache.Get(email)
	if found {
		return queries.RequestError(c, http.StatusBadRequest, "A password reset has already been sent to this email address. Please check your email for the link.")
	}

	token := utils.GenerateSecretCode(3)
	infrastructures.Cache.Set(email, token, time.Minute*10)

	// Send email
	go func() {
		err := utils.SendEmail(email, "Password Reset", "Your password reset code is: "+token)
		if err != nil {
			log := logger.CreateLog(fmt.Sprintf("Error on ResetPassword: %s", err.Error()), logger.LogBodyParserError).SetType(logger.LogTypeError).AttachIDs(0, 0, 0)
			_ = log.SendLog()
		}
	}()

	log := logger.CreateLog(fmt.Sprintf("Password reset request for %s", email), logger.LogUserPasswordReset).SetType(logger.LogTypeInfo).AttachIDs(0, 0, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
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
func (u *UserController) ResetPasswordConfirm(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn()

	var body models.PasswordResetConfirm

	if err := c.BodyParser(&body); err != nil {
		return queries.RequestError(c, http.StatusBadRequest, err.Error())
	}

	email := strings.ToLower(body.Email)
	email = strings.TrimSpace(email)

	token, found := infrastructures.Cache.Get(email)
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
		log := logger.CreateLog(fmt.Sprintf("Error on reset password: %s - %s", user.Username, user.Email), logger.LogBadRequest).SetType(logger.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorRequestFailed)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(body.Pass), 10) // Hash password

	user.Password = password
	db.Save(user)

	infrastructures.Cache.Delete(email)

	log := logger.CreateLog(fmt.Sprintf("Password reset for %s", email), logger.LogUserPasswordChanged).SetType(logger.LogTypeInfo).AttachIDs(user.ID, 0, 0)
	_ = log.SendLog()

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
		Success: true,
		Message: "Success updated password",
		Data:    nil,
		Count:   1,
	})
}

// PUT

// UpdateUserByID function
// @Description Update a user by ID
// @Summary updates a user by ID
// @Tags User
// @Produce json
// @Accept json
// @Param config body models.User true "User"
// @Success 200
// @Router /v1/users/id/{id} [put]
func (u *UserController) UpdateUserByID(c *fiber.Ctx) error {
	db := infrastructures.GetDBConn() // DB Conn

	// Params
	id := c.Params("id")

	user := new(models.User)

	if err := db.First(&user, id).Error; err != nil {
		return queries.RequestError(c, http.StatusInternalServerError, err.Error())
	}

	if res := UpdateUser(c, user); !res.Success {
		return queries.RequestError(c, http.StatusInternalServerError, res.Message)
	}

	return c.Status(http.StatusOK).JSON(viewmodels.ResponseHTTP{
		Success: true,
		Message: "Success update user by ID",
		Data:    *user,
		Count:   1,
	})
}

// UpdateUser function
func UpdateUser(c *fiber.Ctx, user *models.User) *viewmodels.ResponseHTTP {
	db := infrastructures.GetDBConn()

	email, password, permissions := user.Email, user.Password, user.Permissions

	res := new(viewmodels.ResponseHTTP)

	if err := c.BodyParser(&user); err != nil {
		res.GenerateError(err.Error())
		return res
	}

	if user.Email != email || !bytes.Equal(user.Password, password) || user.Permissions != permissions {
		res.GenerateError(utils.ErrorBreak)
		return res
	}

	db.Save(user)

	res.GenerateSuccess("Success update user", nil, 0)
	return res
}
