package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/internal"
)

func registerRoutes(router *fiber.Router) {
	userController := internal.GetServiceContainer().GetUser()
	authController := internal.GetServiceContainer().GetAuth()
	jwtController := internal.GetServiceContainer().GetJwt()
	oauthController := internal.GetServiceContainer().GetOAuth()
	deckController := internal.GetServiceContainer().GetDeck()

	(*router).Add("GET", "/user/me", jwtController.IsConnectedMiddleware(domain.PermissionUser), userController.GetMe)
	(*router).Add("GET", "/user/:uuid", userController.GetName)

	(*router).Add("POST", "/security/login", jwtController.IsConnectedMiddleware(domain.PermissionNone), authController.Login)
	(*router).Add("POST", "/security/register", jwtController.IsConnectedMiddleware(domain.PermissionNone), authController.Register)
	(*router).Add("POST", "/security/logout", jwtController.IsConnectedMiddleware(domain.PermissionUser), authController.Logout)
	(*router).Add("POST", "/security/refresh", jwtController.IsConnectedMiddleware(domain.PermissionUser), authController.RefreshToken)

	(*router).Add("GET", "/security/github", oauthController.GithubLogin)
	(*router).Add("GET", "/security/github_callback", oauthController.GithubCallback)

	(*router).Add("GET", "/security/discord", oauthController.DiscordLogin)
	(*router).Add("GET", "/security/discord_callback", oauthController.DiscordCallback)

	(*router).Add("GET", "/security/key", authController.SearchKey)

	(*router).Add("GET", "/deck/owned", jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.GetOwned)
	(*router).Add("GET", "/deck/learning", jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.GetLearning)
	(*router).Add("GET", "/deck/:id", jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.GetByID)
	(*router).Add("POST", "/deck", jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.Create)
}
