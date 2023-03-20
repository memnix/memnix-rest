package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/memnix/memnix-rest/pkg/env"
	"github.com/memnix/memnix-rest/pkg/json"
)

// JSONHelper is the helper for JSON operations
var JSONHelper = json.NewJSON(&json.GoJSON{})

// EnvHelper is the helper for environment variables
var EnvHelper = env.NewMyEnv(&env.OsEnv{})

const (
	ExpirationTimeInHours = 24 // Expiration time in hours
	SQLMaxOpenConns       = 50 // Max number of open connections to the database
	SQLMaxIdleConns       = 10 // Max number of connections in the idle connection pool

	BCryptCost = 11 // Cost to use for BCrypt

	CleaningInterval = 15 * time.Minute // Interval for cleaning the cache

	OauthStateLength   = 16               // Length of the state for oauth
	OauthStateDuration = 10 * time.Minute // Duration of the state for oauth

	RedisMinIdleConns      = 200               // Min number of idle connections in the pool
	RedisPoolSize          = 12000             // Max number of connections in the pool
	RedisPoolTimeout       = 240 * time.Second // Timeout for getting a connection from the pool
	RedisDefaultExpireTime = 6 * time.Hour     // Default expiration time for redis keys
	RedisOwnedExpireTime   = 2 * time.Hour     // Expiration time for owned keys
)

// JwtSigningMethod is the signing method for JWT
var JwtSigningMethod = jwt.SigningMethodHS256

// PasswordConfigStruct is the struct for the password config
type PasswordConfigStruct struct {
	Iterations uint32 // Number of iterations to use for Argon2ID
	Memory     uint32 // Memory to use for Argon2ID
	Threads    uint8  // Number of threads to use for Argon2ID
	KeyLen     uint32 // Key length to use for Argon2ID
	SaltLen    uint32 // Salt length to use for Argon2ID
}

// IsProduction returns true if the app is in production
func IsProduction() bool {
	return EnvHelper.GetEnv("APP_ENV") != "dev"
}

// IsDevelopment returns true if the app is in development
func IsDevelopment() bool {
	return EnvHelper.GetEnv("APP_ENV") == "dev"
}
