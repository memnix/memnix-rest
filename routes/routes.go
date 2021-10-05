package routes

import (
	"memnixrest/handlers"

	"github.com/gofiber/fiber/v2"
)

// New
func New() *fiber.App {
	// Create new app
	app := fiber.New()

	// Api group
	api := app.Group("/api")

	// v1 group "/api/v1"
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		return c.Next()
	})

	api.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	// Users
	v1.Get("/users", handlers.GetAllUsers)          // Get all users
	v1.Get("/user/id/:id", handlers.GetUserByID)    // Get user by id
	v1.Post("/user/new", handlers.CreateNewUser)    // Create a new user
	v1.Put("/user/id/:id", handlers.UpdateUserByID) // Update an user using his id

	return app

}
