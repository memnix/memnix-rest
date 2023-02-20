package config

import (
	"github.com/golang-jwt/jwt"
	"github.com/memnix/memnix-rest/pkg/env"
	"github.com/memnix/memnix-rest/pkg/json"
)

var JSONHelper = json.NewJSON(&json.GoJson{})

var EnvHelper = env.NewMyEnv(&env.OsEnv{})

const (
	Argon2IDThreads  = 1 // Number of threads to use for Argon2ID
	Argon2IDTime     = 1 // Number of iterations to use for Argon2ID
	Argon2IDMem      = 16 * MB
	Argon2IDKey      = 16 // Key length to use for Argon2ID
	Argon2IDSalt     = 16 // Salt length to use for Argon2ID
	BCryptCost       = 12 // Cost to use for BCrypt
	Argon2HashLength = 6  // Argon2HashLength is the argon2 hash length

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

// PasswordConfig is the default config for the password hash (Argon2ID)
var PasswordConfig = PasswordConfigStruct{
	Iterations: Argon2IDTime,
	Memory:     Argon2IDMem,
	Threads:    Argon2IDThreads,
	KeyLen:     Argon2IDKey,
	SaltLen:    Argon2IDSalt,
}
