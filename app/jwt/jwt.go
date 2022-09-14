package jwt

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/pkg/database"
	"os"
	"strings"
)

var SecretKey string   // SecretKey env variable
var AuthDebugMode bool // AuthDebugMode env variable

func Init() {
	SecretKey = os.Getenv("SECRET")                   // SecretKey env variable
	AuthDebugMode = os.Getenv("AUTH_DEBUG") == "true" // AuthDebugMode env variable
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
func IsConnected(c *fiber.Ctx) (int, models.ResponseAuth) {
	db := database.DBConn          // DB Conn
	tokenString := extractToken(c) // Extract token
	var user models.User           // User object

	// Parse token
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		// Return error
		return fiber.StatusForbidden, models.ResponseAuth{
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
		log := models.CreateLog(fmt.Sprintf("Error on check auth: %s", res.Error), models.LogLoginError).SetType(models.LogTypeError).AttachIDs(user.ID, 0, 0)
		_ = log.SendLog()                         // Send log
		c.Status(fiber.StatusInternalServerError) // InternalServerError Status
		// return error
		return fiber.StatusInternalServerError, models.ResponseAuth{
			Success: false,
			Message: "Failed to get the user. Try to logout/login. Otherwise, contact the support",
			User:    user,
		}
	}

	// User is connected
	return fiber.StatusOK, models.ResponseAuth{
		Success: true,
		Message: "User is connected",
		User:    user,
	}
}
