package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal"
)

func registerRoutes(router *fiber.Router) {
	baristaController := internal.GetServiceContainer().GetBarista()
	klientoController := internal.GetServiceContainer().GetKliento()
	(*router).Add("GET", "/barista", baristaController.GetName)

	(*router).Add("GET", "/kliento", klientoController.GetName)
	(*router).Add("GET", "/kliento/:name", klientoController.SetName)
}
