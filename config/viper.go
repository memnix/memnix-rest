package config

import (
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfigStruct
	Database DatabaseConfigStruct
	Redis    RedisConfigStruct
	Log      LogConfigStruct
	Auth     AuthConfigStruct
	Tracing  TracingConfigStruct
}

type TracingConfigStruct struct {
	URL string
}

// ServerConfigStruct is the server configuration.
type ServerConfigStruct struct {
	Port        string
	AppVersion  string
	JaegerURL   string
	Host        string
	FrontendURL string
}

// DatabaseConfigStruct is the database configuration.
type DatabaseConfigStruct struct {
	DSN string
}

// RedisConfigStruct is the redis configuration.
type RedisConfigStruct struct {
	Addr         string
	Password     string
	MinIdleConns int
	PoolSize     int
	PoolTimeout  int
}

// LogConfigStruct is the log configuration.
type LogConfigStruct struct {
	Level string
}

// AuthConfigStruct is the auth configuration.
type AuthConfigStruct struct {
	JWTSecret     string
	JWTHeaderLen  int
	JWTExpiration int
	Discord       oauth.DiscordConfig
	Github        oauth.GithubConfig
	Bcryptcost    int
}

// UseConfig loads a file from given path
func UseConfig(filename string) (*Config, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
