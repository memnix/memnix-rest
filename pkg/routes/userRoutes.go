package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerUserRoutes(r fiber.Router) {
	// Get
	r.Get("/users", controllers.GetAllUsers)        // Get all users
	r.Get("/users/id/:id", controllers.GetUserByID) // Get user by ID

	r.Get("/secret/migrate", controllers.SecretMigrate)

	// Put
	r.Put("/users/id/:id", controllers.UpdateUserByID) // Update an user using his ID

	//TODO: User Management
}
