package queries

import (
	"github.com/gofiber/fiber/v2"
	"memnixrest/app/models"
	"net/http"
)

func AuthError(c *fiber.Ctx, auth models.ResponseAuth) error {
	return c.Status(http.StatusUnauthorized).JSON(models.ResponseHTTP{
		Success: false,
		Message: auth.Message,
		Data:    nil,
		Count:   0,
	})
}

func RequestError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(models.ResponseHTTP{
		Success: false,
		Message: message,
		Data:    nil,
		Count:   0,
	})
}
