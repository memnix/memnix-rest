package infrastructures

import (
	"github.com/memnix/memnix-rest/config"
)

var AppConfig Config

func GetAppConfig() Config {
	return AppConfig
}

type Config struct {
	GithubConfig
	DiscordConfig
}

func InitOauth() {
	AppConfig.GithubConfig.ClientID = config.EnvHelper.GetEnv("GITHUB_CLIENT_ID")
	AppConfig.GithubConfig.ClientSecret = config.EnvHelper.GetEnv("GITHUB_CLIENT_SECRET")

	AppConfig.DiscordConfig.ClientID = config.EnvHelper.GetEnv("DISCORD_CLIENT_ID")
	AppConfig.DiscordConfig.ClientSecret = config.EnvHelper.GetEnv("DISCORD_CLIENT_SECRET")
	AppConfig.DiscordConfig.URL = config.EnvHelper.GetEnv("DISCORD_URL")
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
}

type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	URL          string
}
