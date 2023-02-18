package config

import "time"

const (
	DiodeLoggerSize = 1000
	DiodeLoggerTime = 10 * time.Millisecond

	MaxBackupLogFiles = 5
	MaxSizeLogFiles   = 50 // megabytes

	RedisMinIdleConns = 200
	RedisPoolSize     = 12000
	RedisPoolTimeout  = 240 * time.Second
	RedisHost         = ":6379"

	MB = 1024 * 1024
)
