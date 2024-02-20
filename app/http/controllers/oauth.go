package controllers

import (
	"fmt"
	views2 "github.com/memnix/memnix-rest/app/http/httpViews"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/memnix/memnix-rest/pkg/random"
	"go.opentelemetry.io/otel/attribute"
)

// OAuthController is the controller for the OAuth routes.
type OAuthController struct {
	auth auth.IUseCase // auth usecase
	auth.IAuthRedisRepository
}

// NewOAuthController creates a new OAuthController.
func NewOAuthController(auth auth.IUseCase, redisRepository auth.IAuthRedisRepository) OAuthController {
	return OAuthController{auth: auth, IAuthRedisRepository: redisRepository}
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
	state, _ := random.GetRandomGeneratorInstance().GenerateSecretCode(config.OauthStateLength)
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s",
		oauth.GetGithubClientID(),
		oauth.GetCallbackURL()+"/v2/security/github_callback",
		state,
	)
	// Save the state in the cache
	if err := a.IAuthRedisRepository.SetState(c.UserContext(), state); err != nil {
		return err
	}
	if err := c.Redirect(redirectURL, fiber.StatusSeeOther); err != nil {
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
	if ok, _ := a.IAuthRedisRepository.HasState(c.UserContext(), state); !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	// get the access token from github
	accessToken, err := oauth.GetGithubAccessToken(c.UserContext(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	// get the user from github
	user, err := oauth.GetGithubData(c.UserContext(), accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	var githubUser domain.GithubLogin
	err = GetJSONHelperInstance().GetJSONHelper().Unmarshal([]byte(user), &githubUser)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	// log the user
	jwtToken, err := a.auth.LoginOauth(c.UserContext(), githubUser.ToUser())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	// Delete the state from the cache
	if err = a.IAuthRedisRepository.DeleteState(c.UserContext(), state); err != nil {
		log.WithContext(c.UserContext()).Error("failed to delete state from cache", slog.Any("error", err))
	}

	return c.Redirect(oauth.GetFrontendURL()+"/callback/"+jwtToken, fiber.StatusSeeOther)
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
	state, _ := random.GetRandomGeneratorInstance().GenerateSecretCode(config.OauthStateLength)
	if err := a.IAuthRedisRepository.SetState(c.UserContext(), state); err != nil {
		return err
	}

	redirectURL := oauth.GetDiscordURL() + "&state=" + state

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
	_, span := infrastructures.GetTracerInstance().Start(c.UserContext(), "DiscordCallback")
	defer span.End()
	// get the code from the query string
	code := c.Query("code")
	state := c.Query("state")

	span.SetAttributes(attribute.String("code", code), attribute.String("state", state))
	if ok, _ := a.IAuthRedisRepository.HasState(c.UserContext(), state); !ok {
		log.WithContext(c.UserContext()).Warn("state not found", slog.String("state", state))
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	// get the access token from discord
	accessToken, err := oauth.GetDiscordAccessToken(c.UserContext(), code)
	if err != nil {
		log.WithContext(c.UserContext()).Error("failed to get access token from discord", slog.Any("error", err))
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	if accessToken == "" {
		log.WithContext(c.UserContext()).Error("access token is empty")
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	// get the user from discord
	user, err := oauth.GetDiscordData(c.UserContext(), accessToken)
	if err != nil {
		log.WithContext(c.UserContext()).Error("failed to get user from discord", slog.Any("error", err))
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	var discordUser domain.DiscordLogin
	// print the user to the console
	err = GetJSONHelperInstance().GetJSONHelper().Unmarshal([]byte(user), &discordUser)
	if err != nil {
		log.WithContext(c.UserContext()).Error("failed to unmarshal discord user", slog.Any("error", err))
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	if discordUser == (domain.DiscordLogin{}) {
		log.WithContext(c.UserContext()).Error("discord user is empty")
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	// log the user
	jwtToken, err := a.auth.LoginOauth(c.UserContext(), discordUser.ToUser())
	if err != nil {
		log.WithContext(c.UserContext()).Error("failed to login user", slog.Any("error", err))
		return c.Status(fiber.StatusUnauthorized).JSON(views2.NewLoginTokenVM("", views2.InvalidCredentials))
	}

	if err = a.IAuthRedisRepository.DeleteState(c.UserContext(), state); err != nil {
		log.WithContext(c.UserContext()).Error("error deleting state", slog.Any("error", err))
	}

	return c.Redirect(oauth.GetFrontendURL()+"/callback/"+jwtToken, fiber.StatusSeeOther)
}
