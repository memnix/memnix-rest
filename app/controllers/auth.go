package controllers

import (
	"fmt"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/app/queries"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = os.Getenv("SECRET")

// Register method to create a new user
// @Description Register a new user
// @Summary registers a new user
// @Tags Auth
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Param username body string true "Username"
// @Success 200 {object} models.User
// @Failure 404 "Error"
// @Failure 403 "Forbidden"
// @Router /register [post]
func Register(c *fiber.Ctx) error {
	var data map[string]string
	db := database.DBConn // DB Conn

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if len(data["password"]) > utils.MaxPasswordLen || len(data["username"]) > utils.MaxUsernameLen || len(data["email"]) > utils.MaxEmailLen {
		log := models.CreateLog(fmt.Sprintf("Error on register: %s - %s", data["username"], data["email"]), models.LogBadRequest).SetType(models.LogTypeWarning).AttachIDs(0, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorRequestFailed)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Username: data["username"],
		Email:    strings.ToLower(data["email"]),
		Password: password,
	}

	if err := db.Create(&user).Error; err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on register: %s - %s", data["username"], data["email"]), models.LogAlreadyUsedEmail).SetType(models.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorAlreadyUsedEmail)
	}

	log := models.CreateLog(fmt.Sprintf("Register: %s - %s", user.Username, user.Email), models.LogUserRegister).SetType(models.LogTypeInfo).AttachIDs(user.ID, 0, 0)
	_ = log.SendLog()

	return c.JSON(user)
}

// Login method to login user and return access with fresh token
// @Description Login user and return access with fresh token
// @Summary logins user and return access with fresh token
// @Tags Auth
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} models.User
// @Failure 404 "Error"
// @Failure 400 "Incorrect password"
// @Failure 500 "Internal error"
// @Router /login [post]
func Login(c *fiber.Ctx) error {
	var data map[string]string
	db := database.DBConn // DB Conn

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	db.Where("email = ?", strings.ToLower(data["email"])).First(&user)

	// handle error
	if user.ID == 0 { //default Id when return nil
		log := models.CreateLog(fmt.Sprintf("Error on login: %s", data["email"]), models.LogIncorrectEmail).SetType(models.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect email or password !",
		})
	}

	// match password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		log := models.CreateLog(fmt.Sprintf("Error on login: %s", data["email"]), models.LogIncorrectPassword).SetType(models.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		return c.JSON(fiber.Map{
			"message": "Incorrect email or password !",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 336).Unix(), //14 day
	}) // expires after 2 weeks

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		log := models.CreateLog(fmt.Sprintf("Error on login: %s", err.Error()), models.LogLoginError).SetType(models.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error when logging in !",
		})
	}

	log := models.CreateLog(fmt.Sprintf("Login: %s - %s", user.Username, user.Email), models.LogUserLogin).SetType(models.LogTypeInfo).AttachIDs(user.ID, 0, 0)
	if err = log.SendLog(); err != nil {
		fmt.Println(err)
	}

	return c.JSON(fiber.Map{
		"message": "Login Succeeded",
		"token":   token,
	})
}

func extractToken(c *fiber.Ctx) string {
	token := c.Get("Authorization")
	// Normally Authorization HTTP header.
	onlyToken := strings.Split(token, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

// User method to get connected user
// @Description To get connected user
// @Summary  gets connected user
// @Tags Auth
// @Produce json
// @Success 200 {object} models.User
// @Failure 404 "Error"
// @Failure 401 "Forbidden"
// @Security ApiKeyAuth
// @Router /user [get]
func User(c *fiber.Ctx) error {
	db := database.DBConn // DB Conn

	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	var user models.User

	db.Where("id = ?", claims["Issuer"]).First(&user)

	return c.JSON(user)
}

func AuthDebugMode(c *fiber.Ctx) models.ResponseAuth {
	db := database.DBConn // DB Conn
	var user models.User

	if res := db.Where("id = ?", 6).First(&user); res.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return models.ResponseAuth{
			Success: false,
			Message: "Failed to get the user. Try to logout/login. Otherwise, contact the support",
		}
	}

	return models.ResponseAuth{
		Success: true,
		Message: "Authenticated",
		User:    user,
	}
}

func CheckAuth(c *fiber.Ctx, p models.Permission) models.ResponseAuth {
	db := database.DBConn // DB Conn
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)

		return models.ResponseAuth{
			Message: "Unauthenticated",
			Success: false,
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	var user models.User

	if res := db.Where("id = ?", claims["Issuer"]).First(&user); res.Error != nil {
		log := models.CreateLog(fmt.Sprintf("Error on check auth: %s", err.Error()), models.LogLoginError).SetType(models.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		c.Status(fiber.StatusInternalServerError)
		return models.ResponseAuth{
			Success: false,
			Message: "Failed to get the user. Try to logout/login. Otherwise, contact the support",
		}
	}

	if user.Permissions < p {
		log := models.CreateLog(fmt.Sprintf("Permission error: %s | had %s but tried %s", user.Email, user.Permissions.ToString(), p.ToString()), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()
		c.Status(fiber.StatusUnauthorized)
		return models.ResponseAuth{
			Success: false,
			Message: "You don't have the right permissions to perform this request.",
		}
	}

	return models.ResponseAuth{
		Success: true,
		Message: "Authenticated",
		User:    user,
	}
}

// Logout method to de-auth connected user and delete token
// @Description Logout to de-auth connected user and delete token
// @Summary logouts and de-auth connected user and delete token
// @Tags Auth
// @Produce json
// @Success 200 "Logout"
// @Failure 404 "Error"
// @Failure 401 "Forbidden"
// @Security ApiKeyAuth[user]
// @Router /logout [post]
func Logout(c *fiber.Ctx) error {
	auth := CheckAuth(c, models.PermUser) // Check auth
	if !auth.Success {
		return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
			Success: false,
			Message: auth.Message,
			Data:    nil,
			Count:   0,
		})
	}

	log := models.CreateLog(fmt.Sprintf("Logout: %s - %s", auth.User.Username, auth.User.Email), models.LogUserLogout).SetType(models.LogTypeInfo).AttachIDs(auth.User.ID, 0, 0)
	_ = log.SendLog()

	return c.JSON(fiber.Map{
		"message": "successfully logged out !",
		"token":   "",
	})
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(SecretKey), nil
}
