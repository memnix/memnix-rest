package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerDeckRoutes(r fiber.Router) {

	// Get
	r.Get("/decks", controllers.GetAllDecks)                 // Get all decks
	r.Get("/decks/public", controllers.GetAllPublicDecks)    // Get all public decks
	r.Get("/decks/user/:userID", controllers.GetAllSubDecks) // Get all decks the user is sub to
	r.Get("/decks/id/:id", controllers.GetDeckByID)          // Get deck by ID

	// Post
	r.Post("/decks/new", controllers.CreateNewDeck)                            // Create a new deck
	r.Post("/decks/:deckID/subscribe", controllers.SubToDeck)     // Subscribe to a deck
	r.Post("/decks/:deckID/unsubscribe", controllers.UnSubToDeck) // Unsubscribe to a deck
}
