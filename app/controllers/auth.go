package controllers

import (
	"fmt"
	"github.com/memnix/memnixrest/app/auth"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/logger"
	"github.com/memnix/memnixrest/pkg/models"
	"github.com/memnix/memnixrest/pkg/queries"
	"github.com/memnix/memnixrest/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Register function to create a new user
// @Description Create a new user
// @Summary creates a new user
// @Tags Auth
// @Produce json
// @Param credentials body models.RegisterStruct true "Credentials"
// @Success 200 {object} models.User
// @Failure 403 "Forbidden"
// @Router /v1/register [post]
func Register(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var data models.RegisterStruct // Data object

	if err := c.BodyParser(&data); err != nil {
		return err
	} // Parse body

	// Register checks
	if len(data.Password) > utils.MaxPasswordLen || len(data.Username) > utils.MaxUsernameLen || len(data.Email) > utils.MaxEmailLen {
		log := logger.CreateLog(fmt.Sprintf("Error on register: %s - %s", data.Username, data.Email), logger.LogBadRequest).SetType(logger.LogTypeWarning).AttachIDs(0, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorRequestFailed)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 10) // Hash password
	user := models.User{
		Username: data.Username,
		Email:    strings.ToLower(data.Email),
		Password: password,
	} // Create object

	//TODO: manual checking for unique username and email
	if err := db.Create(&user).Error; err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on register: %s - %s", data.Username, data.Email), logger.LogAlreadyUsedEmail).SetType(logger.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorAlreadyUsedEmail)
	} // Add user to DB

	// Create log
	log := logger.CreateLog(fmt.Sprintf("Register: %s - %s", user.Username, user.Email), logger.LogUserRegister).SetType(logger.LogTypeInfo).AttachIDs(user.ID, 0, 0)
	_ = log.SendLog() // Send log

	return c.JSON(user) // Return user
}

// Login function to log in a user and return access with fresh token
// @Description Login the user and return a fresh token
// @Summary logins user and return a fresh token
// @Tags Auth
// @Produce json
// @Param credentials body models.LoginStruct true "Credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 "Incorrect password or email"
// @Failure 500 "Internal error"
// @Router /v1/login [post]
func Login(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	var data models.LoginStruct // Data object

	if err := c.BodyParser(&data); err != nil {
		return err
	} // Parse body

	var user models.User // User object

	db.Where("email = ?", strings.ToLower(data.Email)).First(&user) // Get user

	// handle error
	if user.ID == 0 { // default ID when return nil
		// Create log
		log := logger.CreateLog(fmt.Sprintf("Error on login: %s", data.Email), logger.LogIncorrectEmail).SetType(logger.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()                // Send log
		c.Status(fiber.StatusBadRequest) // BadRequest Status
		// return error message as Json object
		return c.JSON(models.LoginResponse{
			Message: "Incorrect email or password !",
			Token:   "",
		})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password)); err != nil {
		c.Status(fiber.StatusBadRequest) // BadRequest Status
		log := logger.CreateLog(fmt.Sprintf("Error on login: %s", data.Email), logger.LogIncorrectPassword).SetType(logger.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog() // Send log
		// return error message as Json object
		return c.JSON(models.LoginResponse{
			Message: "Incorrect email or password !",
			Token:   "",
		})
	}

	// Create token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 336).Unix(), // 14 day
	}) // expires after 2 weeks

	token, err := claims.SignedString([]byte(auth.SecretKey)) // Sign token
	if err != nil {
		log := logger.CreateLog(fmt.Sprintf("Error on login: %s", err.Error()), logger.LogLoginError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()                         // Send log
		c.Status(fiber.StatusInternalServerError) // InternalServerError Status
		// return error message as Json object
		return c.JSON(models.LoginResponse{
			Message: "Incorrect email or password !",
			Token:   "",
		})
	}

	log := logger.CreateLog(fmt.Sprintf("Login: %s - %s", user.Username, user.Email), logger.LogUserLogin).SetType(logger.LogTypeInfo).AttachIDs(user.ID, 0, 0)
	_ = log.SendLog() // Send log

	// return token as Json object
	return c.JSON(models.LoginResponse{
		Message: "Login Succeeded",
		Token:   token,
	})
}

// User function to get connected user
// @Description Get connected user
// @Summary  gets connected user
// @Tags Auth
// @Produce json
// @Success 200 {object} models.ResponseAuth
// @Failure 401 "Forbidden"
// @Security Beaver
// @Router /v1/user [get]
func User(c *fiber.Ctx) error {
	// statusCode, response := IsConnected(c) // Check if connected

	user := new(models.PublicUser)

	localUser, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	user.Set(&localUser) // Set user

	responseUser := models.ResponsePublicAuth{
		Success: true,
		Message: "User is connected",
		User:    *user,
	}

	return c.Status(fiber.StatusOK).JSON(responseUser) // Return response
}

// Logout function to log user logout
// @Description Logout the user and create a record in the log
// @Summary logouts the user
// @Tags Auth
// @Produce json
// @Success 200 "Success"
// @Failure 401 "Forbidden"
// @Security Beaver
// @Router /v1/logout [post]
func Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return queries.RequestError(c, http.StatusUnauthorized, utils.ErrorForbidden)
	}

	// Create log
	log := logger.CreateLog(fmt.Sprintf("Logout: %s - %s", user.Username, user.Email), logger.LogUserLogout).SetType(logger.LogTypeInfo).AttachIDs(user.ID, 0, 0)
	_ = log.SendLog()

	// Return response with success
	return c.JSON(fiber.Map{
		"message": "successfully logged out !",
		"token":   "",
	})
}
