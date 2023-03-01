package config

import "time"

const (
	DiodeLoggerSize = 1000
	DiodeLoggerTime = 10 * time.Millisecond

	MaxBackupLogFiles = 5
	MaxSizeLogFiles   = 50 // megabytes
	LogChannelSize    = 200

	RedisMinIdleConns = 200
	RedisPoolSize     = 12000
	RedisPoolTimeout  = 240 * time.Second
	RedisHost         = "redis:6379"

	MB = 1024 * 1024

	Base10  = 10 // Base10 is the base 10
	BitSize = 32 // BitSize is the bit size

	JwtTokenHeaderLen = 2 // JwtTokenHeaderLen is the jwt token header length

	MaxPasswordLength = 72 // Max password length based on bcrypt limit
	MinPasswordLength = 8  // Min password length
)
