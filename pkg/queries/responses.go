package queries

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/viewmodels"
	"net/http"
)

func AuthError(c *fiber.Ctx, auth *viewmodels.ResponseAuth) error {
	return c.Status(http.StatusUnauthorized).JSON(viewmodels.ResponseHTTP{
		Success: false,
		Message: auth.Message,
		Data:    nil,
		Count:   0,
	})
}

func RequestError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(viewmodels.ResponseHTTP{
		Success: false,
		Message: message,
		Data:    nil,
		Count:   0,
	})
}
