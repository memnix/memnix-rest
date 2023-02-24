package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/internal"
)

func registerRoutes(router *fiber.Router) {
	userController := internal.GetServiceContainer().GetUser()
	klientoController := internal.GetServiceContainer().GetKliento()
	authController := internal.GetServiceContainer().GetAuth()
	jwtController := internal.GetServiceContainer().GetJwt()
	oauthController := internal.GetServiceContainer().GetOAuth()

	(*router).Add("GET", "/user/:uuid", userController.GetName)

	(*router).Add("GET", "/kliento", klientoController.GetName)
	(*router).Add("GET", "/kliento/:name", klientoController.SetName)

	(*router).Add("POST", "/security/login", jwtController.IsConnectedMiddleware(domain.PermissionNone), authController.Login)
	(*router).Add("POST", "/security/register", jwtController.IsConnectedMiddleware(domain.PermissionNone), authController.Register)
	(*router).Add("POST", "/security/logout", jwtController.IsConnectedMiddleware(domain.PermissionUser), authController.Logout)
	(*router).Add("POST", "/security/refresh", jwtController.IsConnectedMiddleware(domain.PermissionUser), authController.RefreshToken)

	(*router).Add("GET", "/security/github", oauthController.GithubLogin)
	(*router).Add("GET", "/security/github_callback", oauthController.GithubCallback)
}
