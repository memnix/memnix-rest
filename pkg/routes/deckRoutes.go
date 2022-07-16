package routes

import (
	"github.com/memnix/memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerDeckRoutes(r fiber.Router) { // Get
	r.Get("/decks", controllers.GetAllDecks)                    // Get all decks
	r.Get("/decks/public", controllers.GetAllPublicDecks)       // Get all public decks
	r.Get("/decks/available", controllers.GetAllAvailableDecks) // Get all available decks
	r.Get("/decks/editor", controllers.GetAllEditorDecks)       // Get all decks the user is editor
	r.Get("/decks/sub", controllers.GetAllSubDecks)             // Get all decks the user is sub to
	r.Get("/decks/:deckID", controllers.GetDeckByID)            // Get deck by ID
	r.Get("/decks/:deckID/users", controllers.GetAllSubUsers)   // Get all sub users

	// Post
	r.Post("/decks/new", controllers.CreateNewDeck)                             // Create a new deck
	r.Post("/decks/:deckID/subscribe", controllers.SubToDeck)                   // Subscribe to a deck
	r.Post("/decks/:deckID/unsubscribe", controllers.UnSubToDeck)               // Unsubscribe to a deck
	r.Post("/decks/private/:key/:code/subscribe", controllers.SubToPrivateDeck) // Subscribe to a private deck using key and code
	r.Post("/decks/:deckID/publish", controllers.PublishDeckRequest)            // Request to publish a deck

	// Put
	r.Put("/decks/:deckID/edit", controllers.UpdateDeckByID) // Update a deck by ID

	// Delete
	r.Delete("/decks/:deckID", controllers.DeleteDeckById) // Delete a deck by ID
}
