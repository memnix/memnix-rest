package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/memnix/memnix-rest/pkg/env"
	"github.com/memnix/memnix-rest/pkg/json"
)

var JSONHelper = json.NewJSON(&json.GoJson{})

var EnvHelper = env.NewMyEnv(&env.OsEnv{})

const (
	ExpirationTimeInHours = 24 // Expiration time in hours
	SQLMaxOpenConns       = 50 // Max number of open connections to the database
	SQLMaxIdleConns       = 10 // Max number of connections in the idle connection pool

	BCryptCost = 11 // Cost to use for BCrypt

	CleaningInterval = 15 * time.Minute

	OauthStateLength   = 16               // Length of the state for oauth
	OauthStateDuration = 10 * time.Minute // Duration of the state for oauth

	RedisMinIdleConns      = 200
	RedisPoolSize          = 12000
	RedisPoolTimeout       = 240 * time.Second
	RedisDefaultExpireTime = 6 * time.Hour
	RedisOwnedExpireTime   = 2 * time.Hour
)

var JwtSigningMethod = jwt.SigningMethodHS256 // JWTSigningMethod is the signing method for JWT

// PasswordConfigStruct is the struct for the password config
type PasswordConfigStruct struct {
	Iterations uint32 // Number of iterations to use for Argon2ID
	Memory     uint32 // Memory to use for Argon2ID
	Threads    uint8  // Number of threads to use for Argon2ID
	KeyLen     uint32 // Key length to use for Argon2ID
	SaltLen    uint32 // Salt length to use for Argon2ID
}

func IsProduction() bool {
	return EnvHelper.GetEnv("APP_ENV") != "dev"
}

func IsDevelopment() bool {
	return EnvHelper.GetEnv("APP_ENV") == "dev"
}
