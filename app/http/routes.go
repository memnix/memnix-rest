package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal"
)

func registerRoutes(router *fiber.Router) {
	userController := internal.GetServiceContainer().GetUser()
	klientoController := internal.GetServiceContainer().GetKliento()
	authController := internal.GetServiceContainer().GetAuth()

	(*router).Add("GET", "/user/:uuid", userController.GetName)

	(*router).Add("GET", "/kliento", klientoController.GetName)
	(*router).Add("GET", "/kliento/:name", klientoController.SetName)

	(*router).Add("POST", "/auth/login", authController.Login)
	(*router).Add("POST", "/auth/register", authController.Register)
	(*router).Add("POST", "/auth/logout", authController.Logout)
	(*router).Add("POST", "/auth/refresh", authController.RefreshToken)
}
