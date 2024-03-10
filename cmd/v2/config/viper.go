package config

import (
	"log/slog"

	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/spf13/viper"
)

// Config holds the configuration for the application.
type Config struct {
	Server    ServerConfig
	Log       LogConfig
	Pgx       infrastructures.PgxConfig
	Auth      AuthConfig
	Sentry    infrastructures.SentryConfig
	Redis     infrastructures.RedisConfig
	Ristretto infrastructures.RistrettoConfig
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
	Discord       oauth.DiscordConfig
	Github        oauth.GithubConfig
	JWTSecret     string
	JWTHeaderLen  int
	JWTExpiration int
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
