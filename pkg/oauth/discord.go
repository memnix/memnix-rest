package oauth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/rs/zerolog/log"
)

// GetDiscordAccessToken gets the access token from Discord
func GetDiscordAccessToken(code string) (string, error) {
	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s&scope=identify,email",
		infrastructures.AppConfig.DiscordConfig.ClientID,
		infrastructures.AppConfig.DiscordConfig.ClientSecret,
		config.GetCurrentURL()+"/v2/security/discord_callback",
		code,
	)))

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://discord.com/api/oauth2/token",
		reqBody,
	)
	if reqerr != nil {
		log.Info().Msg("Request failed")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Info().Msg("Response failed")
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	// Represents the response received from Github
	type discordAccessTokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		Expires      int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp discordAccessTokenResponse
	err := config.JSONHelper.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

// GetDiscordData gets the user data from Discord
func GetDiscordData(accessToken string) (string, error) {
	req, reqerr := http.NewRequest(
		"GET",
		"https://discord.com/api/users/@me",
		nil,
	)
	if reqerr != nil {
		log.Info().Msg("Request failed")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Info().Msg("Response failed")
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	return string(respbody), nil
}
