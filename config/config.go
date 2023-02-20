package config

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/memnix/memnix-rest/pkg/env"
	"github.com/memnix/memnix-rest/pkg/json"
)

var JSONHelper = json.NewJSON(&json.GoJson{})

var EnvHelper = env.NewMyEnv(&env.OsEnv{})

const (
	ExpirationTimeInHours = 24 // Expiration time in hours
)

var (
	JwtSigningMethod = jwt.SigningMethodHS256 // JWTSigningMethod is the signing method for JWT
)

// PasswordConfigStruct is the struct for the password config
type PasswordConfigStruct struct {
	Iterations uint32 // Number of iterations to use for Argon2ID
	Memory     uint32 // Memory to use for Argon2ID
	Threads    uint8  // Number of threads to use for Argon2ID
	KeyLen     uint32 // Key length to use for Argon2ID
	SaltLen    uint32 // Salt length to use for Argon2ID
}
