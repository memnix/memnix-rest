package oauth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"
)

const (
	discordTokenURL = "https://discord.com/api/oauth2/token" //nolint:gosec //This is a URL, not a password.
	discordAPIURL   = "https://discord.com/api/users/@me"    //nolint:gosec //This is a URL, not a password.
)

// Represents the response received from Discord.
type discordAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	Expires      int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// GetDiscordAccessToken gets the access token from Discord.
func GetDiscordAccessToken(ctx context.Context, code string) (string, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(
		"client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s&scope=identify,email",
		discordConfig.ClientID,
		discordConfig.ClientSecret,
		GetCallbackURL()+"/v2/security/discord_callback",
		code,
	))

	// POST request to set URL.
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		discordTokenURL,
		reqBody,
	)
	if err != nil {
		log.WithContext(ctx).Error("Failed to get Discord access token", slog.Any("error", err))
		return "", errors.Wrap(err, RequestFailed)
	}

	if req == nil || req.Body == nil || req.Header == nil {
		log.WithContext(ctx).Error("Failed to get Discord access token", slog.Any("error", err))
		return "", errors.New(RequestFailed)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Get the response.
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.WithContext(ctx).Error("Failed to get Discord access token", slog.Any("error", resperr))
		return "", errors.Wrap(resperr, ResponseFailed)
	}

	defer func(resp *http.Response) {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}(resp)

	// Response body converted to stringified JSON.
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithContext(ctx).Error("failed to read resp.body", slog.Any("error", err))
		return "", err
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse.
	var ghresp discordAccessTokenResponse
	err = GetJSONHelperInstance().GetJSONHelper().Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us).
	return ghresp.AccessToken, nil
}

// GetDiscordData gets the user data from Discord.
func GetDiscordData(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		discordAPIURL,
		nil,
	)
	if err != nil {
		return "", errors.Wrap(err, RequestFailed)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Get the response.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, ResponseFailed)
	}

	defer func(resp *http.Response) {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}(resp)
	// Response body converted to stringified JSON.
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respbody), nil
}
