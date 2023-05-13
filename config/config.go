package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memnix/memnix-rest/pkg/env"
	"github.com/memnix/memnix-rest/pkg/json"
)

// JSONHelper is the helper for JSON operations
var JSONHelper = json.NewJSON(&json.SonicJSON{})

// EnvHelper is the helper for environment variables
var EnvHelper = env.NewMyEnv(&env.OsEnv{})

const (
	ExpirationTimeInHours = 24 // ExpirationTimeInHours is the expiration time for the JWT token
	SQLMaxOpenConns       = 50 // SQLMaxOpenConns is the max number of connections in the open connection pool
	SQLMaxIdleConns       = 10 // SQLMaxIdleConns is the max number of connections in the idle connection pool

	BCryptCost = 11 // BCryptCost is the cost for bcrypt

	OauthStateLength   = 16               // OauthStateLength is the length of the state for oauth
	OauthStateDuration = 10 * time.Minute // OauthStateDuration is the duration for the state for oauth

	RedisMinIdleConns      = 200               // RedisMinIdleConns is the minimum number of idle connections in the pool
	RedisPoolSize          = 12000             // RedisPoolSize is the maximum number of connections allocated by the pool at a given time
	RedisPoolTimeout       = 240 * time.Second // RedisPoolTimeout is the amount of time a connection can be used before being closed
	RedisDefaultExpireTime = 6 * time.Hour     // RedisDefaultExpireTime is the default expiration time for keys
	RedisOwnedExpireTime   = 2 * time.Hour     // RedisOwnedExpireTime is the expiration time for owned keys

	CacheExpireTime = 10 * time.Second // CacheExpireTime is the expiration time for the cache
	InfluxDBFreq    = 10 * time.Second // InfluxDBFreq is the frequency for writing to InfluxDB

	DeckSecretCodeLength = 10  //DeckSecretCodeLength is the length of the secret code for decks
	GCThresholdPercent   = 0.7 // GCThresholdPercent is the threshold for garbage collection
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
