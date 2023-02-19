package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal"
)

func registerRoutes(router *fiber.Router) {
	userController := internal.GetServiceContainer().GetUser()
	klientoController := internal.GetServiceContainer().GetKliento()
	(*router).Add("GET", "/user/:uuid", userController.GetName)

	(*router).Add("GET", "/kliento", klientoController.GetName)
	(*router).Add("GET", "/kliento/:name", klientoController.SetName)
}
