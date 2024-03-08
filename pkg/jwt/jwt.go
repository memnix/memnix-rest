package jwt

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

type InstanceSingleton struct {
	jwtInstance Instance
}

var (
	jwtInstance *InstanceSingleton //nolint:gochecknoglobals //Singleton
	jwtOnce     sync.Once          //nolint:gochecknoglobals //Singleton
)

func GetJwtInstance() *InstanceSingleton {
	jwtOnce.Do(func() {
		jwtInstance = &InstanceSingleton{}
	})
	return jwtInstance
}

func (j *InstanceSingleton) GetJwt() Instance {
	return j.jwtInstance
}

func (j *InstanceSingleton) SetJwt(instance Instance) {
	j.jwtInstance = instance
}

type Instance struct {
	SigningMethod         jwt.SigningMethod
	PublicKey             ed25519.PublicKey
	PrivateKey            ed25519.PrivateKey
	HeaderLen             int
	ExpirationTimeInHours int
}

// NewJWTInstance return a new JwtInstance with the given parameters.
func NewJWTInstance(headerLen, expirationTime int,
	publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey,
) Instance {
	return Instance{
		HeaderLen:             headerLen,
		PublicKey:             publicKey,
		PrivateKey:            privateKey,
		SigningMethod:         jwt.SigningMethodEdDSA,
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
func (instance Instance) GenerateToken(_ context.Context, userID int32) (string, error) {
	// Create the Claims for the token
	claims := jwt.NewWithClaims(instance.SigningMethod, jwt.RegisteredClaims{
		Issuer:    utils.ConvertInt32ToStr(userID),    // Issuer is the user id
		ExpiresAt: instance.CalculateExpirationTime(), // ExpiresAt is the expiration time
	})

	// Sign and get the complete encoded token as a string using the secret
	token, err := claims.SignedString(instance.PrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign")
	}

	return token, nil
}

// VerifyToken verifies a jwt token
// and returns the user id and an error.
func (Instance) VerifyToken(token *jwt.Token) (int32, error) {
	// claims is of type jwt.MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {
		// Get the issuer from the claims and convert it to uint
		userID, err := utils.ConvertStrToInt32(claims["iss"].(string))
		if err != nil {
			return 0, err
		}

		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// GetToken gets a jwt.Token token from a string
// and returns the jwt.Token and an error.
func (instance Instance) GetToken(_ context.Context, token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key.
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return instance.PublicKey, nil // Return the secret key as the signing key
	})
}

func (Instance) GetExpirationTime(token *jwt.Token) int64 {
	// Safe type assertion
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0
	}
	return int64(claims["exp"].(float64))
}

// ExtractToken function to extract token from header.
func (instance Instance) ExtractToken(token string) string {
	// Normally Authorization HTTP header.
	onlyToken := strings.Split(token, " ") // Split token.
	if len(onlyToken) == instance.HeaderLen {
		return onlyToken[1] // Return only token.
	}
	return "" // Return empty string.
}

// GetConnectedUserID gets the user id from a jwt token.
func (instance Instance) GetConnectedUserID(ctx context.Context, tokenHeader string) (int32, error) {
	if tokenHeader == "" {
		return 0, errors.New("empty token")
	}

	// Get the token from the Authorization header.
	tokenString := instance.ExtractToken(tokenHeader)

	if tokenString == "" {
		return 0, errors.New("empty token")
	}

	token, err := instance.GetToken(ctx, tokenString)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	// Check if the token is valid.
	userID, err := instance.VerifyToken(token)
	if err != nil {
		return 0, err
	}

	if userID == 0 {
		return 0, errors.New("invalid token")
	}

	return userID, nil
}

// CalculateExpirationTime returns the expiration time.
func (instance Instance) CalculateExpirationTime() *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(instance.ExpirationTimeInHours)))
}
