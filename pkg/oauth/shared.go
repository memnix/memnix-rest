package oauth

import (
	"sync"

	"github.com/memnix/memnix-rest/pkg/json"
)

const (
	// RequestFailed is the error message when a request fails.
	RequestFailed = "request failed"
	// ResponseFailed is the error message when a response fails.
	ResponseFailed = "response failed"
)

// GlobalConfig is the global configuration for oauth.
type GlobalConfig struct {
	CallbackURL string // CallbackURL is the callback URL for oauth.
	FrontendURL string // FrontendURL is the frontend URL for oauth.
}

// JSONHelperSingleton is a singleton for the json helper.
type JSONHelperSingleton struct {
	jsonHelper json.Helper
}

var (
	instance *JSONHelperSingleton //nolint:gochecknoglobals //Singleton
	once     sync.Once            //nolint:gochecknoglobals //Singleton
)

// GetJSONHelperInstance gets the json helper instance.
func GetJSONHelperInstance() *JSONHelperSingleton {
	once.Do(func() {
		instance = &JSONHelperSingleton{
			jsonHelper: json.NewJSON(&json.NativeJSON{}),
		}
	})
	return instance
}

// GetJSONHelper gets the json helper.
func (j *JSONHelperSingleton) GetJSONHelper() json.Helper {
	return j.jsonHelper
}

// SetJSONHelper sets the json helper.
func (j *JSONHelperSingleton) SetJSONHelper(helper json.Helper) {
	j.jsonHelper = json.NewJSON(helper)
}

var (
	githubConfig  GithubConfig  //nolint:gochecknoglobals //Viper will handle the global configuration.
	discordConfig DiscordConfig //nolint:gochecknoglobals //Viper will handle the global configuration.
	oauthConfig   GlobalConfig  //nolint:gochecknoglobals //Viper will handle the global configuration.
)

type GithubConfig struct {
	ClientID     string // ClientID is the client ID for Github.
	ClientSecret string // ClientSecret is the client secret for Github.
}

type DiscordConfig struct {
	ClientID     string // ClientID is the client ID for Discord.
	ClientSecret string // ClientSecret is the client secret for Discord.
	URL          string // URL is the URL for Discord.
}

// GetGithubClientID gets the Github client ID.
func GetGithubClientID() string {
	return githubConfig.ClientID
}

// GetCallbackURL gets the callback URL.
func GetCallbackURL() string {
	return oauthConfig.CallbackURL
}

// GetFrontendURL gets the frontend URL.
func GetFrontendURL() string {
	return oauthConfig.FrontendURL
}

// GetDiscordURL gets the Discord client ID.
func GetDiscordURL() string {
	return discordConfig.URL
}

// InitGithub initializes the Github oauth configuration.
func InitGithub(config GithubConfig) {
	githubConfig = config
}

// InitDiscord initializes the Discord oauth configuration.
func InitDiscord(config DiscordConfig) {
	discordConfig = config
}

// SetOauthConfig sets the oauth configuration.
func SetOauthConfig(config GlobalConfig) {
	oauthConfig = config
}
