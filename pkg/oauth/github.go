package oauth

import (
	"bytes"
	"context"
	"fmt"
	"github.com/memnix/memnix-rest/views"
	"io"
	"net/http"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetGithubAccessToken gets the access token from Github using the code
func GetGithubAccessToken(code string) (string, error) {
	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     infrastructures.GetAppConfig().GithubConfig.ClientID,
		"client_secret": infrastructures.GetAppConfig().GithubConfig.ClientSecret,
		"code":          code,
	}
	requestJSON, _ := config.JSONHelper.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil || req == nil || req.Body == nil || req.Header == nil {
		log.Debug().Err(reqerr).Msg("github.go: GetGithubAccessToken: Request failed (reqerr != nil || req == nil || req.Body == nil || req.Header == nil)")
		return "", errors.Wrap(reqerr, views.RequestFailed)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil || resp == nil || resp.Body == nil {
		log.Debug().Err(resperr).Msg("github.go: GetGithubAccessToken: Response failed (resperr != nil || resp == nil || resp.Body == nil)")
		return "", errors.Wrap(reqerr, views.ResponseFailed)
	}

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	err := config.JSONHelper.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

// GetGithubData gets the user data from Github using the access token
func GetGithubData(accessToken string) (string, error) {
	// Get request to a set URL
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		log.Info().Msg(views.RequestFailed)
		return "", err
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Info().Msg(views.ResponseFailed)
		return "", err
	}

	// Read the response as a byte slice
	respbody, _ := io.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody), nil
}
