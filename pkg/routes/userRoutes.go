package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerUserRoutes(r fiber.Router) {
	// Get
	r.Get("/users", controllers.GetAllUsers)                           // Get all users
	r.Get("/users/id/:id", controllers.GetUserByID)                    // Get user by ID
	r.Get("/users/discord/:discordID", controllers.GetUserByDiscordID) // Get user by discordID

	// Post
	r.Post("/users/new", controllers.CreateNewUser) // Create a new user

	// Put
	r.Put("/users/id/:id", controllers.UpdateUserByID) // Update an user using his ID
}
