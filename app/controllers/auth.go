package controllers

import (
	"memnixrest/app/models"
	"memnixrest/app/queries"
	"memnixrest/pkg/database"
	"memnixrest/pkg/utils"
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

func Register(c *fiber.Ctx) error {
	var data map[string]string
	db := database.DBConn // DB Conn

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Username: data["username"],
		Email:    strings.ToLower(data["email"]),
		Password: password,
	}

	if err := db.Create(&user).Error; err != nil {
		return queries.RequestError(c, http.StatusForbidden, utils.ErrorAlreadyUsedEmail)
	}

	log := queries.CreateLog(models.LogUserRegister, "Register: "+user.Username)
	_ = queries.CreateUserLog(user.ID, log)

	return c.JSON(user)
}

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
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	// match password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password!",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 168).Unix(), //7 day
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "error when logging in !",
		})
	}

	cookie := fiber.Cookie{
		Name:     "memnix-jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 168),
		HTTPOnly: true,
		SameSite: "none",
		Secure:   true,
	}
	c.Cookie(&cookie)

	log := queries.CreateLog(models.LogUserLogin, "Login: "+user.Username)
	_ = queries.CreateUserLog(user.ID, log)

	return c.JSON(fiber.Map{
		"message": "Login Succeeded",
		//"token": token,
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("memnix-jwt")
	db := database.DBConn // DB Conn

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	db.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func CheckAuth(c *fiber.Ctx, p models.Permission) models.ResponseAuth {
	cookie := c.Cookies("memnix-jwt")
	db := database.DBConn // DB Conn
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)

		return models.ResponseAuth{
			Message: "Unauthenticated",
			Success: false,
		}
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	if res := db.Where("id = ?", claims.Issuer).First(&user); res.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return models.ResponseAuth{
			Success: false,
			Message: "Failed to get the user. Try to logout/login. Otherwise, contact the support",
		}
	}

	if user.Permissions < p {
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

	cookie := fiber.Cookie{
		Name:     "memnix-jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "none",
		Secure:   true,
	}
	c.Cookie(&cookie)

	log := queries.CreateLog(models.LogUserLogout, "Logout: "+auth.User.Username)
	_ = queries.CreateUserLog(auth.User.ID, log)

	return c.JSON(fiber.Map{
		"message": "successfully logged out !",
	})
}
