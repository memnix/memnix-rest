package config

import "time"

const (
	DiodeLoggerSize = 1000
	DiodeLoggerTime = 10 * time.Millisecond

	MaxBackupLogFiles = 5
	MaxSizeLogFiles   = 20 // megabytes
	LogChannelSize    = 200

	MB = 1024 * 1024

	Base10  = 10 // Base10 is the base 10
	BitSize = 32 // BitSize is the bit size

	JwtTokenHeaderLen = 2 // JwtTokenHeaderLen is the jwt token header length

	MaxPasswordLength = 72 // Max password length based on bcrypt limit
	MinPasswordLength = 8  // Min password length

	Localhost = "http://localhost:1815"
	ApiHost   = "https://beta.memnix.app"

	FrontHost      = "https://memnix.corentings.dev"
	FrontHostLocal = "http://localhost:3000"
)

func GetCurrentURL() string {
	if IsProduction() {
		return ApiHost
	}

	return Localhost
}

func GetFrontURL() string {
	if IsProduction() {
		return FrontHost
	}

	return FrontHostLocal
}
