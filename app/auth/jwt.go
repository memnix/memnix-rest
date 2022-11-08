package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/pkg/logger"
	"github.com/memnix/memnixrest/viewmodels"
	"os"
	"strings"
)

var SecretKey string // SecretKey env variable
var _ bool           // AuthDebugMode env variable

func Init() {
	SecretKey = os.Getenv("SECRET")       // SecretKey env variable
	_ = os.Getenv("AUTH_DEBUG") == "true" // AuthDebugMode env variable
}

// extractToken function to extract token from header
func extractToken(c *fiber.Ctx) string {
	token := c.Get("Authorization") // Get token from header
	// Normally Authorization HTTP header.
	onlyToken := strings.Split(token, " ") // Split token
	if len(onlyToken) == 2 {
		return onlyToken[1] // Return only token
	}
	return "" // Return empty string
}

// jwtKeyFunc function to get the key for the token
func jwtKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(SecretKey), nil // Return secret key
}

// IsConnected function to check if user is connected
func IsConnected(c *fiber.Ctx) (int, viewmodels.ResponseAuth) {
	db := infrastructures.GetDBConn() // DB Conn
	tokenString := extractToken(c)    // Extract token
	var user models.User              // User object

	// Parse token
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		// Return error
		return fiber.StatusForbidden, viewmodels.ResponseAuth{
			Success: false,
			Message: "Failed to get the user. Try to logout/login. Otherwise, contact the support",
			User:    user,
		}
	}
	// Check if token is valid
	claims := token.Claims.(jwt.MapClaims)

	// Get user from token
	if res := db.Where("id = ?", claims["iss"]).First(&user); res.Error != nil {
		// Generate log
		log := logger.CreateLog(fmt.Sprintf("Error on check auth: %s", res.Error), logger.LogLoginError).SetType(logger.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()                         // Send log
		c.Status(fiber.StatusInternalServerError) // InternalServerError Status
		// return error
		return fiber.StatusInternalServerError, viewmodels.ResponseAuth{
			Success: false,
			Message: "Failed to get the user. Try to logout/login. Otherwise, contact the support",
			User:    user,
		}
	}

	// User is connected
	return fiber.StatusOK, viewmodels.ResponseAuth{
		Success: true,
		Message: "User is connected",
		User:    user,
	}
}
