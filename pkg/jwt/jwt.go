package jwt

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/utils"
)

// GenerateToken generates a jwt token from a user id
// and returns the token and an error
//
// It's signing method is defined in utils.JwtSigningMethod
// It's expiration time is defined in utils.GetExpirationTime
// It's secret key is defined in the environment variable SECRET_KEY
// see: utils/config.go for more information
func GenerateToken(ctx context.Context, userID uint) (string, error) {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "GenerateToken")
	defer span.End()
	// Create the Claims for the token
	claims := jwt.NewWithClaims(config.JwtSigningMethod, jwt.RegisteredClaims{
		Issuer:    utils.ConvertUIntToStr(userID), // Issuer is the user id
		ExpiresAt: utils.GetExpirationTime(),      // ExpiresAt is the expiration time
	})

	// Sign and get the complete encoded token as a string using the secret
	token, err := claims.SignedString([]byte(utils.GetSecretKey()))
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyToken verifies a jwt token
// and returns the user id and an error
func VerifyToken(token *jwt.Token) (uint, error) {
	// claims is of type jwt.MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {
		// Get the issuer from the claims and convert it to uint
		userID, err := utils.ConvertStrToUInt(claims["iss"].(string))
		if err != nil {
			return 0, err
		}

		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// GetToken gets a jwt.Token token from a string
// and returns the jwt.Token and an error
func GetToken(token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key.
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(utils.GetSecretKey()), nil // Return the secret key as the signing key
	})
}

func GetExpirationTime(token *jwt.Token) int64 {
	claims := token.Claims.(jwt.MapClaims)
	return int64(claims["exp"].(float64))
}

// extractToken function to extract token from header
func extractToken(token string) string {
	// Normally Authorization HTTP header.
	onlyToken := strings.Split(token, " ") // Split token
	if len(onlyToken) == config.JwtTokenHeaderLen {
		return onlyToken[1] // Return only token
	}
	return "" // Return empty string
}

// GetConnectedUserID gets the user id from a jwt token
func GetConnectedUserID(tokenHeader string) (uint, error) {
	// Get the token from the Authorization header
	tokenString := extractToken(tokenHeader)

	token, err := GetToken(tokenString)
	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	userID, err := VerifyToken(token)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
