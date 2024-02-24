package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/services"
)

func registerRoutes(router *fiber.Router) {
	serviceContainer := services.DefaultServiceContainer()
	userController := serviceContainer.User()
	authController := serviceContainer.Auth()
	jwtController := serviceContainer.Jwt()
	oauthController := serviceContainer.OAuth()
	deckController := serviceContainer.Deck()

	(*router).Get("/health", func(c *fiber.Ctx) error { return c.Status(http.StatusOK).SendString("ok") })

	(*router).Get("/user/me",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), userController.GetMe)

	(*router).Get("/user/:uuid",
		userController.GetName)

	(*router).Post("/security/login",
		jwtController.IsConnectedMiddleware(domain.PermissionNone), authController.Login)

	(*router).Post("/security/register",
		jwtController.IsConnectedMiddleware(domain.PermissionNone), authController.Register)

	(*router).Post("/security/logout",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), authController.Logout)

	(*router).Post("/security/refresh",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), authController.RefreshToken)

	(*router).Get("/security/github",
		oauthController.GithubLogin)

	(*router).Get("/security/github_callback",
		oauthController.GithubCallback)

	(*router).Get("/security/discord",
		oauthController.DiscordLogin)

	(*router).Get("/security/discord_callback",
		oauthController.DiscordCallback)

	(*router).Get("/deck/owned",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.GetOwned)

	(*router).Get("/deck/public",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.GetPublic)

	(*router).Get("/deck/learning",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.GetLearning)

	(*router).Get("/deck/:id",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.GetByID)

	(*router).Post("/deck",
		jwtController.IsConnectedMiddleware(domain.PermissionUser), deckController.Create)
}
