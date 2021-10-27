package controllers

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string
	db := database.DBConn // DB Conn

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: password,
	}

	db.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	db := database.DBConn // DB Conn

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	db.Where("email = ?", data["email"]).First(&user)

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
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
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
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

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

func CheckAuth(c *fiber.Ctx) models.ResponseAuth {
	cookie := c.Cookies("memnix-jwt")
	db := database.DBConn // DB Conn
	response := new(models.ResponseAuth)
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		response.Message = "Unauthentified"
		response.Success = false
		return *response
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	db.Where("id = ?", claims.Issuer).First(&user)

	response.User = user
	response.Success = true
	response.Message = "Authentified"

	return *response
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "memnix-jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "successfully logged out !",
	})
}
