package config

import (
	"os"
	"time"

	"github.com/memnix/memnix-rest/pkg/crypto"
	"github.com/memnix/memnix-rest/pkg/json"
	myJwt "github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"golang.org/x/crypto/ed25519"
)

// JSONHelper is the helper for JSON operations
var JSONHelper = json.NewJSON(&json.SonicJSON{})

const (
	ExpirationTimeInHours = 24 // ExpirationTimeInHours is the expiration time for the JWT token
	SQLMaxOpenConns       = 10 // SQLMaxOpenConns is the max number of connections in the open connection pool
	SQLMaxIdleConns       = 1  // SQLMaxIdleConns is the max number of connections in the idle connection pool

	BCryptCost = 11 // BCryptCost is the cost for bcrypt

	OauthStateLength = 16 // OauthStateLength is the length of the state for oauth

	RedisDefaultExpireTime = 6 * time.Hour // RedisDefaultExpireTime is the default expiration time for keys

	CacheExpireTime = 10 * time.Second // CacheExpireTime is the expiration time for the cache

	GCThresholdPercent = 0.7 // GCThresholdPercent is the threshold for garbage collection

	GCLimit = 1024 * 1024 * 1024 // GCLimit is the limit for garbage collection

	RistrettoMaxCost     = 5 * MB // RistrettoMaxCost is the maximum cost
	RistrettoBufferItems = 32     // RistrettoBufferItems is the number of items per get buffer
	RistrettoNumCounters = 1e4    // RistrettoNumCounters is the number of counters

	MB = 1024 * 1024 // MB is the number of bytes in a megabyte

	MaxPasswordLength = 72 // MaxPasswordLength is the max password length
	MinPasswordLength = 8  // MinPasswordLength is the min password length

	SentryFlushTimeout = 2 * time.Second // SentryFlushTimeout is the timeout for flushing sentry
)

var JwtInstance myJwt.Instance

func GetJwtInstance() myJwt.Instance {
	return JwtInstance
}

// PasswordConfigStruct is the struct for the password config
type PasswordConfigStruct struct {
	Iterations uint32 // Iterations to use for Argon2ID
	Memory     uint32 // Memory to use for Argon2ID
	Threads    uint8  // Threads to use for Argon2ID
	KeyLen     uint32 // KeyLen to use for Argon2ID
	SaltLen    uint32 // SaltLen to use for Argon2ID
}

var (
	ed25519PrivateKey = ed25519.PrivateKey{}
	ed25519PublicKey  = ed25519.PublicKey{}
)

func ParseEd25519Key() error {
	publicKey, privateKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return err
	}

	ed25519PrivateKey = privateKey
	ed25519PublicKey = publicKey

	otelzap.L().Info("✅ Created ed25519 keys")

	return nil
}

// GetEd25519PrivateKey returns the ed25519 private key
func GetEd25519PrivateKey() ed25519.PrivateKey {
	return ed25519PrivateKey
}

// GetEd25519PublicKey returns the ed25519 public key
func GetEd25519PublicKey() ed25519.PublicKey {
	return ed25519PublicKey
}

func GetConfigPath() string {
	if IsDevelopment() {
		return "./config/config-local"
	}

	return "./config/config-prod"
}

func IsProduction() bool {
	return os.Getenv("APP_ENV") != "dev"
}

func IsDevelopment() bool {
	return os.Getenv("APP_ENV") == "dev"
}

func GetCallbackURL() string {
	return os.Getenv("CALLBACK_URL")
}
