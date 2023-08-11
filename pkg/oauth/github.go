package oauth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/views"
	"github.com/pkg/errors"
)

// Represents the response received from Github
type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

const (
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
	githubAPIURL         = "https://api.github.com/user"
)

// GetGithubAccessToken gets the access token from Github using the code
func GetGithubAccessToken(ctx context.Context, code string) (string, error) {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "GetGithubAccessToken")
	defer span.End()
	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     infrastructures.GetAppConfig().GithubConfig.ClientID,
		"client_secret": infrastructures.GetAppConfig().GithubConfig.ClientSecret,
		"code":          code,
	}
	requestJSON, _ := config.JSONHelper.Marshal(requestBodyMap)

	// POST request to set URL
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		githubAccessTokenURL,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil || req == nil || req.Body == nil || req.Header == nil {
		return "", errors.Wrap(err, views.RequestFailed)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, views.ResponseFailed)
	}

	defer resp.Body.Close()

	// Response body converted to stringified JSON
	respbody, _ := io.ReadAll(resp.Body)

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	err = config.JSONHelper.Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

// GetGithubData gets the user data from Github using the access token
func GetGithubData(ctx context.Context, accessToken string) (string, error) {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "GetGithubData")
	defer span.End()
	// Get request to a set URL
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return "", err
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Read the response as a byte slice
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert byte slice to string and return
	return string(respbody), nil
}
