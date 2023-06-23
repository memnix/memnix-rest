package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/memnix/memnix-rest/pkg/env"
	"github.com/memnix/memnix-rest/pkg/json"
	"golang.org/x/crypto/ed25519"
)

// JSONHelper is the helper for JSON operations
var JSONHelper = json.NewJSON(&json.SonicJSON{})

// EnvHelper is the helper for environment variables
var EnvHelper = env.NewMyEnv(&env.OsEnv{})

const (
	ExpirationTimeInHours = 24 // ExpirationTimeInHours is the expiration time for the JWT token
	SQLMaxOpenConns       = 20 // SQLMaxOpenConns is the max number of connections in the open connection pool
	SQLMaxIdleConns       = 5  // SQLMaxIdleConns is the max number of connections in the idle connection pool

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

	DeckSecretCodeLength = 10  // DeckSecretCodeLength is the length of the secret code for decks
	GCThresholdPercent   = 0.7 // GCThresholdPercent is the threshold for garbage collection

	GCLimit = 1024 * 1024 * 1024 // GCLimit is the limit for garbage collection

	GormPrometheusRefreshInterval = 15 // GormPrometheusRefreshInterval is the refresh interval for gorm prometheus

	RistrettoMaxCost     = 5 * MB // RistrettoMaxCost is the maximum cost
	RistrettoBufferItems = 32     // RistrettoBufferItems is the number of items per get buffer
	RistrettoNumCounters = 1e4    // RistrettoNumCounters is the number of counters
)

// JwtSigningMethod is the signing method for JWT
var JwtSigningMethod = jwt.SigningMethodEdDSA

// PasswordConfigStruct is the struct for the password config
type PasswordConfigStruct struct {
	Iterations uint32 // Iterations to use for Argon2ID
	Memory     uint32 // Memory to use for Argon2ID
	Threads    uint8  // Threads to use for Argon2ID
	KeyLen     uint32 // KeyLen to use for Argon2ID
	SaltLen    uint32 // SaltLen to use for Argon2ID
}

// IsProduction returns true if the app is in production
func IsProduction() bool {
	return EnvHelper.GetEnv("APP_ENV") != "dev"
}

// IsDevelopment returns true if the app is in development
func IsDevelopment() bool {
	return EnvHelper.GetEnv("APP_ENV") == "dev"
}

var (
	ed25519PrivateKey = ed25519.PrivateKey{}
	ed25519PublicKey  = ed25519.PublicKey{}
)

// GetEd25519PrivateKey returns the ed25519 private key
func GetEd25519PrivateKey() ed25519.PrivateKey {
	return ed25519PrivateKey
}

// GetEd25519PublicKey returns the ed25519 public key
func GetEd25519PublicKey() ed25519.PublicKey {
	return ed25519PublicKey
}

// ParseEd25519PrivateKey parses the ed25519 private key
func ParseEd25519PrivateKey() error {
	key, err := os.ReadFile("./config/keys/ed25519_private.pem")
	if err != nil {
		return err
	}

	privateKey, err := jwt.ParseEdPrivateKeyFromPEM(key)
	if err != nil {
		return err
	}

	ed25519PrivateKey = privateKey.(ed25519.PrivateKey)
	return nil
}

// ParseEd25519PublicKey parses the ed25519 public key
func ParseEd25519PublicKey() error {
	key, err := os.ReadFile("./config/keys/ed25519_public.pem")
	if err != nil {
		return err
	}

	publicKey, err := jwt.ParseEdPublicKeyFromPEM(key)
	if err != nil {
		return err
	}

	ed25519PublicKey = publicKey.(ed25519.PublicKey)
	return nil
}
