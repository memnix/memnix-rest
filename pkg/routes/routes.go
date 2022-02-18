package routes

import (
	"memnixrest/app/controllers"

	_ "memnixrest/docs"

	"github.com/gofiber/fiber/v2/middleware/cors"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	// Create new app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost, *",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", swagger.Handler) // default

	// Api group
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	// Auth
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/user", controllers.User)
	api.Post("/logout", controllers.Logout)

	// v1 group "/api/v1"
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		return c.Next()
	})

	v1.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	// Register routes
	registerUserRoutes(v1) // /v1/users/
	registerDeckRoutes(v1) // /v1/decks/
	registerCardRoutes(v1) // /v1/cards/

	return app
}
