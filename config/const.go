package config

import "time"

const (
	DiodeLoggerSize = 1000                  // DiodeLoggerSize is the size of the diode logger
	DiodeLoggerTime = 10 * time.Millisecond // DiodeLoggerTime is the time of the diode logger

	MaxBackupLogFiles = 5   // MaxBackupLogFiles is the max number of backup log files
	MaxSizeLogFiles   = 20  // MaxSizeLogFiles is the max size of log files in MB
	LogChannelSize    = 100 // LogChannelSize is the size of the log channel

	MB = 1024 * 1024 // MB is the number of bytes in a megabyte

	Base10  = 10 // Base10 is the base 10
	BitSize = 32 // BitSize is the bit size

	JwtTokenHeaderLen = 2 // JwtTokenHeaderLen is the jwt token header length

	MaxPasswordLength = 72 // Max password length based on bcrypt limit
	MinPasswordLength = 8  // Min password length

	Localhost = "http://localhost:1815"   // Localhost is the localhost url
	APIHost   = "https://beta.memnix.app" // APIHost is the api host url

	FrontHost      = "https://memnix.corentings.dev" // FrontHost is the front host url
	FrontHostLocal = "http://localhost:3000"         // FrontHostLocal is the front host url for local development
)

// GetCurrentURL returns the current url
func GetCurrentURL() string {
	if IsProduction() {
		return APIHost
	}

	return Localhost
}

// GetFrontURL returns the front url
func GetFrontURL() string {
	if IsProduction() {
		return FrontHost
	}

	return FrontHostLocal
}
