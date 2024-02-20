package oauth

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	githubAccessTokenURL = "https://github.com/login/oauth/access_token" //nolint:gosec //This is a URL, not a password.
	githubAPIURL         = "https://api.github.com/user"                 //nolint:gosec //This is a URL, not a password.
)

// githubAccessTokenResponse Represents the response received from Github.
type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GetGithubAccessToken gets the access token from Github using the code.
func GetGithubAccessToken(ctx context.Context, code string) (string, error) {
	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     githubConfig.ClientID,
		"client_secret": githubConfig.ClientSecret,
		"code":          code,
	}
	requestJSON, _ := GetJSONHelperInstance().GetJSONHelper().Marshal(requestBodyMap)

	// POST request to set URL
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		githubAccessTokenURL,
		bytes.NewBuffer(requestJSON),
	)
	if err != nil || req == nil || req.Body == nil || req.Header == nil {
		return "", errors.Wrap(err, RequestFailed)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
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
	respbody, _ := io.ReadAll(resp.Body)

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	err = GetJSONHelperInstance().GetJSONHelper().Unmarshal(respbody, &ghresp)
	if err != nil {
		return "", err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

// GetGithubData gets the user data from Github using the access token.
func GetGithubData(ctx context.Context, accessToken string) (string, error) {
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

	defer func(resp *http.Response) {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}(resp)
	// Read the response as a byte slice
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert byte slice to string and return
	return string(respbody), nil
}
