package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/pkg/cacheset"
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/memnix/memnix-rest/pkg/random"
	"github.com/memnix/memnix-rest/views"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// OAuthController is the controller for the OAuth routes
type OAuthController struct {
	auth     auth.IUseCase   // auth usecase
	cacheSet *cacheset.Cache // cache set
}

// NewOAuthController creates a new OAuthController
func NewOAuthController(auth auth.IUseCase, cacheSet *cacheset.Cache) OAuthController {
	return OAuthController{auth: auth, cacheSet: cacheSet}
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
	state, _ := random.GenerateSecretCode(config.OauthStateLength)
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s",
		infrastructures.AppConfig.GithubConfig.ClientID,
		"http://localhost:1815/v2/security/github_callback",
		state,
	)
	// Save the state in the cache
	a.cacheSet.Set(state, config.OauthStateDuration)
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
	state := c.Query("state")

	// check if the state is valid
	if !a.cacheSet.Exists(state) {
		log.Debug().Msg("invalid state")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

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

	// Delete the state from the cache
	a.cacheSet.Delete(state)

	return c.Redirect("http://localhost:3000/callback/"+jwtToken, fiber.StatusSeeOther)
}

// DiscordLogin redirects the user to the discord login page
//
//	@Summary		Redirects the user to the discord login page
//	@Description	Redirects the user to the discord login page
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Success		302	{string}	string					"redirecting to github login"
//	@Failure		500	{object}	views.HTTPResponseVM	"internal server error"
//	@Router			/v2/security/discord [get]
func (a *OAuthController) DiscordLogin(c *fiber.Ctx) error {
	// Create the dynamic redirect URL for login
	state, _ := random.GenerateSecretCode(config.OauthStateLength)
	a.cacheSet.Set(state, config.OauthStateDuration)

	redirectURL := infrastructures.AppConfig.DiscordConfig.URL + "&state=" + state

	err := c.Redirect(redirectURL, fiber.StatusSeeOther)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "redirecting to discord login", "redirect_url": redirectURL})
}

// DiscordCallback handles the callback from discord
//
//	@Summary		Handles the callback from discord
//	@Description	Handles the callback from discord
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string					true	"code from discord"
//	@Success		200		{object}	views.LoginTokenVM		"login token"
//	@Failure		401		{object}	views.HTTPResponseVM	"invalid credentials"
//	@Failure		500		{object}	views.HTTPResponseVM	"internal server error"
//	@Router			/v2/security/discord_callback [get]
func (a *OAuthController) DiscordCallback(c *fiber.Ctx) error {
	// get the code from the query string
	code := c.Query("code")
	state := c.Query("state")

	if !a.cacheSet.Exists(state) {
		log.Debug().Err(errors.New("invalid state")).Msg("invalid state")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	// get the access token from discord
	accessToken, err := oauth.GetDiscordAccessToken(code)
	if err != nil {
		log.Debug().Err(err).Msg("invalid discord access token")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	// get the user from discord
	user, err := oauth.GetDiscordData(accessToken)
	if err != nil {
		log.Debug().Err(err).Msg("invalid discord user")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	var discordUser domain.DiscordLogin
	err = config.JSONHelper.Unmarshal([]byte(user), &discordUser)
	if err != nil {
		log.Debug().Err(err).Msg("can't unmarshal discord user")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	// log the user
	jwtToken, err := a.auth.LoginOauth(discordUser.ToUser())
	if err != nil {
		log.Debug().Err(err).Msg("invalid credentials")
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	a.cacheSet.Delete(state)

	return c.Redirect("http://localhost:3000/callback/"+jwtToken, fiber.StatusSeeOther)
}
