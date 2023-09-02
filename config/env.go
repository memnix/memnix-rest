package config

import "os"

func GetConfigPath() string {
	if IsDevelopment() {
		return "./config/config-local"
	}

	return "./config/config-prod"
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
