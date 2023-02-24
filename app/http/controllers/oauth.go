package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/memnix/memnix-rest/views"
	"github.com/rs/zerolog/log"
)

// OAuthController is the controller for the OAuth routes
type OAuthController struct {
	auth auth.IUseCase // auth usecase
}

// NewOAuthController creates a new OAuthController
func NewOAuthController(auth auth.IUseCase) OAuthController {
	return OAuthController{auth: auth}
}

// GithubLogin redirects the user to the github login page
//
//	@Summary		Redirects the user to the github login page
//	@Description	Redirects the user to the github login page
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Success		302	{string}	string					"redirecting to github login"
//	@Failure		500	{object}	views.HTTPResponseVM	"internal server error"
//	@Router			/v2/security/github [get]
func (a *OAuthController) GithubLogin(c *fiber.Ctx) error {
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		infrastructures.AppConfig.GithubConfig.ClientID,
		"http://localhost:1815/v2/security/github_callback",
	)
	err := c.Redirect(redirectURL, fiber.StatusSeeOther)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "redirecting to github login", "redirect_url": redirectURL})
}

// GithubCallback handles the callback from github
//
//	@Summary		Handles the callback from github
//	@Description	Handles the callback from github
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string					true	"code from github"
//	@Success		200		{object}	views.LoginTokenVM		"login token"
//	@Failure		401		{object}	views.HTTPResponseVM	"invalid credentials"
//	@Failure		500		{object}	views.HTTPResponseVM	"internal server error"
//	@Router			/v2/security/github_callback [get]
func (a *OAuthController) GithubCallback(c *fiber.Ctx) error {
	// get the code from the query string
	code := c.Query("code")

	// get the access token from github
	accessToken, err := oauth.GetGithubAccessToken(code)
	if err != nil {
		log.Debug().Err(err).Msg("invalid github access token")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	// get the user from github
	user, err := oauth.GetGithubData(accessToken)
	if err != nil {
		log.Debug().Err(err).Msg("invalid github user")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	var githubUser domain.GithubLogin
	err = config.JSONHelper.Unmarshal([]byte(user), &githubUser)
	if err != nil {
		log.Debug().Err(err).Msg("can't unmarshal github user")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	// log the user
	jwtToken, err := a.auth.LoginOauth(githubUser.ToUser())
	if err != nil {
		log.Debug().Err(err).Msg("invalid credentials")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewLoginTokenVM(jwtToken, ""))
}
