package config

import (
	"os"
	"time"
)

const (
	GCThresholdPercent = 0.7 // GCThresholdPercent is the threshold for garbage collection

	GCLimit = 1024 * 1024 * 1024 // GCLimit is the limit for garbage collection

	MB = 1024 * 1024 // MB is the number of bytes in a megabyte

	SentryFlushTimeout = 2 * time.Second // SentryFlushTimeout is the timeout for flushing sentry
)

// PasswordConfigStruct is the struct for the password config.
type PasswordConfigStruct struct {
	Iterations uint32 // Iterations to use for Argon2ID
	Memory     uint32 // Memory to use for Argon2ID
	Threads    uint8  // Threads to use for Argon2ID
	KeyLen     uint32 // KeyLen to use for Argon2ID
	SaltLen    uint32 // SaltLen to use for Argon2ID
}

func GetConfigPath() string {
	if IsDevelopment() {
		return "./cmd/v2/config/config-local"
	}

	return "./cmd/v2/config/config-prod"
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
