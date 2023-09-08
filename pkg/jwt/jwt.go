package jwt

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

type Instance struct {
	headerLen             int
	publicKey             ed25519.PublicKey
	privateKey            ed25519.PrivateKey
	signingMethod         jwt.SigningMethod
	ExpirationTimeInHours int
}

// NewJWTInstance return a new JwtInstance with the given parameters
func NewJWTInstance(headerLen, expirationTime int, publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) Instance {
	return Instance{
		headerLen:             headerLen,
		publicKey:             publicKey,
		privateKey:            privateKey,
		signingMethod:         jwt.SigningMethodEdDSA,
		ExpirationTimeInHours: expirationTime,
	}
}

// GenerateToken generates a jwt token from a user id
// and returns the token and an error
//
// It's signing method is defined in utils.JwtSigningMethod
// It's expiration time is defined in utils.GetExpirationTime
// It's secret key is defined in the environment variable SECRET_KEY
// see: utils/config.go for more information
func (instance Instance) GenerateToken(_ context.Context, userID uint) (string, error) {
	// Create the Claims for the token
	claims := jwt.NewWithClaims(instance.signingMethod, jwt.RegisteredClaims{
		Issuer:    utils.ConvertUIntToStr(userID),     // Issuer is the user id
		ExpiresAt: instance.CalculateExpirationTime(), // ExpiresAt is the expiration time
	})

	// Sign and get the complete encoded token as a string using the secret
	token, err := claims.SignedString(instance.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign")
	}

	return token, nil
}

// VerifyToken verifies a jwt token
// and returns the user id and an error
func (Instance) VerifyToken(token *jwt.Token) (uint, error) {
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
func (instance Instance) GetToken(_ context.Context, token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key.
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return instance.publicKey, nil // Return the secret key as the signing key
	})
}

func (Instance) GetExpirationTime(token *jwt.Token) int64 {
	claims := token.Claims.(jwt.MapClaims)
	return int64(claims["exp"].(float64))
}

// extractToken function to extract token from header
func (instance Instance) extractToken(token string) string {
	// Normally Authorization HTTP header.
	onlyToken := strings.Split(token, " ") // Split token
	if len(onlyToken) == instance.headerLen {
		return onlyToken[1] // Return only token
	}
	return "" // Return empty string
}

// GetConnectedUserID gets the user id from a jwt token
func (instance Instance) GetConnectedUserID(ctx context.Context, tokenHeader string) (uint, error) {
	// Get the token from the Authorization header
	tokenString := instance.extractToken(tokenHeader)

	token, err := instance.GetToken(ctx, tokenString)
	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	userID, err := instance.VerifyToken(token)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// CalculateExpirationTime returns the expiration time
func (instance Instance) CalculateExpirationTime() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(instance.ExpirationTimeInHours)))
}
