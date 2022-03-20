package routes

import (
	"github.com/memnix/memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerUserRoutes(r fiber.Router) {
	// Get
	r.Get("/users", controllers.GetAllUsers)        // Get all users
	r.Get("/users/id/:id", controllers.GetUserByID) // Get user by ID

	// Post
	r.Post("/users/settings/:deckID/today", controllers.SetTodayConfig)

	// Put
	r.Put("/users/id/:id", controllers.UpdateUserByID) // Update a user using his ID

	//TODO: User Management
}
