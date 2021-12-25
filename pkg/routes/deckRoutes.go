package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerDeckRoutes(r fiber.Router) {

	// Get
	r.Get("/decks", controllers.GetAllDecks)              // Get all decks
	r.Get("/decks/public", controllers.GetAllPublicDecks) // Get all public decks
	r.Get("/decks/available", controllers.GetAllAvailableDecks)
	r.Get("/decks/sub", controllers.GetAllSubDecks)       // Get all decks the user is sub to
	r.Get("/decks/:id", controllers.GetDeckByID)          // Get deck by ID
	r.Get("/decks/:id/users", controllers.GetAllSubUsers) // Get all sub users

	// Post
	r.Post("/decks/new", controllers.CreateNewDeck)           // Create a new deck
	r.Post("/decks/:id/subscribe", controllers.SubToDeck)     // Subscribe to a deck
	r.Post("/decks/:id/unsubscribe", controllers.UnSubToDeck) // Unsubscribe to a deck

	// Put
	r.Put("/decks/:id/edit", controllers.UpdateDeckByID) // Update a deck by ID

}
