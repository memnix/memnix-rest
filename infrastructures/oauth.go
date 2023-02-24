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
}

func InitOauth() {
	AppConfig.GithubConfig.ClientID = config.EnvHelper.GetEnv("GITHUB_CLIENT_ID")
	AppConfig.GithubConfig.ClientSecret = config.EnvHelper.GetEnv("GITHUB_CLIENT_SECRET")
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
}
