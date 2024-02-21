package config

import (
	"log/slog"

	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/spf13/viper"
)

// Config holds the configuration for the application.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
	Auth     AuthConfig
	Sentry   SentryConfig
}

// SentryConfig holds the configuration for the sentry client.
type SentryConfig struct {
	Debug              bool
	Environment        string
	Release            string
	TracesSampleRate   float64
	ProfilesSampleRate float64
	DSN                string
}

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	Port        string
	AppVersion  string
	JaegerURL   string
	Host        string
	FrontendURL string
	LogLevel    string
}

// DatabaseConfig holds the configuration for the database.
type DatabaseConfig struct {
	DSN string
}

// RedisConfig holds the configuration for the redis client.
type RedisConfig struct {
	Addr         string
	Password     string
	MinIdleConns int
	PoolSize     int
	PoolTimeout  int
}

// LogConfig holds the configuration for the logger.
type LogConfig struct {
	Level string
}

func (logConfig *LogConfig) GetSlogLevel() slog.Level {
	switch logConfig.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// AuthConfig holds the configuration for the authentication.
type AuthConfig struct {
	JWTSecret     string
	JWTHeaderLen  int
	JWTExpiration int
	Discord       oauth.DiscordConfig
	Github        oauth.GithubConfig
	Bcryptcost    int
}

// LoadConfig loads the configuration from a file.
func LoadConfig(filename string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
