package oauth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/pkg/errors"
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
	req, reqerr := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		"https://discord.com/api/oauth2/token",
		reqBody,
	)
	if reqerr != nil || req == nil || req.Body == nil || req.Header == nil {
		log.Debug().Err(reqerr).Msg("discord.go: GetDiscordAccessToken: Request failed (reqerr != nil || req == nil || req.Body == nil || req.Header == nil)")
		return "", errors.Wrap(reqerr, "discord.go: GetDiscordAccessToken: Request failed (reqerr != nil || req == nil || req.Body == nil || req.Header == nil)")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil || resp == nil || resp.Body == nil {
		log.Debug().Err(resperr).Msg("discord.go: GetDiscordAccessToken: Response failed (resperr != nil || resp == nil || resp.Body == nil)")
		return "", errors.Wrap(resperr, "discord.go: GetDiscordAccessToken: Response failed (resperr != nil || resp == nil || resp.Body == nil)")
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
	req, reqerr := http.NewRequestWithContext(context.Background(),
		http.MethodGet,
		"https://discord.com/api/users/@me",
		nil,
	)

	if reqerr != nil || req == nil || req.Body == nil || req.Header == nil {
		log.Debug().Err(reqerr).Msg("discord.go: GetDiscordData: Request failed (reqerr != nil || req == nil || req.Body == nil || req.Header == nil)")
		return "", errors.Wrap(reqerr, "discord.go: GetDiscordData: Request failed (reqerr != nil || req == nil || req.Body == nil || req.Header == nil)")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil || resp == nil || resp.Body == nil {
		log.Debug().Err(resperr).Msg("discord.go: GetDiscordData: Response failed (resperr != nil || resp == nil || resp.Body == nil)")
		return "", errors.Wrap(resperr, "discord.go: GetDiscordData: Response failed (resperr != nil || resp == nil || resp.Body == nil)")
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	return string(respbody), nil
}
