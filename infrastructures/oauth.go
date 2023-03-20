package infrastructures

import (
	"github.com/memnix/memnix-rest/config"
)

// AppConfig is the application configuration.
var AppConfig Config

// GetAppConfig returns the application configuration.
func GetAppConfig() Config {
	return AppConfig
}

// Config is the application configuration.
type Config struct {
	GithubConfig
	DiscordConfig
}

// InitOauth initializes the oauth configuration.
func InitOauth() {
	AppConfig.GithubConfig.ClientID = config.EnvHelper.GetEnv("GITHUB_CLIENT_ID")
	AppConfig.GithubConfig.ClientSecret = config.EnvHelper.GetEnv("GITHUB_CLIENT_SECRET")

	AppConfig.DiscordConfig.ClientID = config.EnvHelper.GetEnv("DISCORD_CLIENT_ID")
	AppConfig.DiscordConfig.ClientSecret = config.EnvHelper.GetEnv("DISCORD_CLIENT_SECRET")
	AppConfig.DiscordConfig.URL = config.EnvHelper.GetEnv("DISCORD_URL")
}

// GithubConfig is the Github oauth configuration.
type GithubConfig struct {
	ClientID     string
	ClientSecret string
}

// DiscordConfig is the Discord oauth configuration.
type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	URL          string
}
