package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/app/auth"
	"github.com/memnix/memnixrest/pkg/models"
	"github.com/memnix/memnixrest/pkg/queries"
	"strings"
)

func IsConnectedMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		path := strings.TrimLeft(c.Path(), "/v1")
		path = strings.TrimRight(path, "/")

		p := routesMap["/"+path].Permission

		if p == models.PermNone {
			return c.Next()
		}

		statusCode, response := auth.IsConnected(c) // Check if connected

		// Check statusCode
		if statusCode != fiber.StatusOK {
			c.Status(statusCode)
			// Return response
			return queries.AuthError(c, &response)
		}

		user := response.User // Get user from response

		// Check permission
		if user.Permissions < p {
			// Log permission error
			log := models.CreateLog(fmt.Sprintf("Permission error: %s | had %s but tried %s", user.Email, user.Permissions.ToString(), p.ToString()), models.LogPermissionForbidden).SetType(models.LogTypeWarning).AttachIDs(user.ID, 0, 0)
			_ = log.SendLog()                  // Send log
			c.Status(fiber.StatusUnauthorized) // Unauthorized Status
			// Return response
			return queries.AuthError(c, &models.ResponseAuth{
				Success: false,
				Message: "You don't have the right permissions to perform this request.",
			})
		}

		// Validate permissions
		c.Locals("user", user) // Set user in locals
		return c.Next()
	}
}
