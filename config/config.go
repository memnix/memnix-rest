package config

import (
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/memnix/memnix-rest/pkg/crypto"
	myJwt "github.com/memnix/memnix-rest/pkg/jwt"
	"golang.org/x/crypto/ed25519"
)

const (
	ExpirationTimeInHours = 24 // ExpirationTimeInHours is the expiration time for the JWT token
	SQLMaxOpenConns       = 10 // SQLMaxOpenConns is the max number of connections in the open connection pool
	SQLMaxIdleConns       = 1  // SQLMaxIdleConns is the max number of connections in the idle connection pool

	OauthStateLength = 16 // OauthStateLength is the length of the state for oauth

	RedisDefaultExpireTime = 6 * time.Hour // RedisDefaultExpireTime is the default expiration time for keys

	CacheExpireTime = 10 * time.Second // CacheExpireTime is the expiration time for the cache

	GCThresholdPercent = 0.7 // GCThresholdPercent is the threshold for garbage collection

	GCLimit = 1024 * 1024 * 1024 // GCLimit is the limit for garbage collection

	RistrettoMaxCost     = 3 * MB // RistrettoMaxCost is the maximum cost
	RistrettoBufferItems = 32     // RistrettoBufferItems is the number of items per get buffer
	RistrettoNumCounters = 1e4    // RistrettoNumCounters is the number of counters

	MB = 1024 * 1024 // MB is the number of bytes in a megabyte

	MaxPasswordLength = 72 // MaxPasswordLength is the max password length
	MinPasswordLength = 8  // MinPasswordLength is the min password length

	SentryFlushTimeout = 2 * time.Second // SentryFlushTimeout is the timeout for flushing sentry
)

type JwtInstanceSingleton struct {
	jwtInstance myJwt.Instance
}

var (
	jwtInstance *JwtInstanceSingleton //nolint:gochecknoglobals //Singleton
	jwtOnce     sync.Once             //nolint:gochecknoglobals //Singleton
)

func GetJwtInstance() *JwtInstanceSingleton {
	jwtOnce.Do(func() {
		jwtInstance = &JwtInstanceSingleton{}
	})
	return jwtInstance
}

func (j *JwtInstanceSingleton) GetJwt() myJwt.Instance {
	return j.jwtInstance
}

func (j *JwtInstanceSingleton) SetJwt(instance myJwt.Instance) {
	j.jwtInstance = instance
}

// PasswordConfigStruct is the struct for the password config.
type PasswordConfigStruct struct {
	Iterations uint32 // Iterations to use for Argon2ID
	Memory     uint32 // Memory to use for Argon2ID
	Threads    uint8  // Threads to use for Argon2ID
	KeyLen     uint32 // KeyLen to use for Argon2ID
	SaltLen    uint32 // SaltLen to use for Argon2ID
}

type KeyManager struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

var (
	keyManagerInstance *KeyManager //nolint:gochecknoglobals //Singleton
	keyManagerOnce     sync.Once   //nolint:gochecknoglobals //Singleton
)

func GetKeyManagerInstance() *KeyManager {
	keyManagerOnce.Do(func() {
		keyManagerInstance = &KeyManager{}
	})
	return keyManagerInstance
}

func (k *KeyManager) GetPrivateKey() ed25519.PrivateKey {
	return k.privateKey
}

func (k *KeyManager) GetPublicKey() ed25519.PublicKey {
	return k.publicKey
}

func (k *KeyManager) ParseEd25519Key() error {
	publicKey, privateKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return err
	}

	k.privateKey = privateKey
	k.publicKey = publicKey

	log.Info("✅ Created ed25519 keys")

	return nil
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
