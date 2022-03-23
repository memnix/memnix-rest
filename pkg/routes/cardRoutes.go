package routes

import (
	"github.com/memnix/memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerCardRoutes(r fiber.Router) {
	// Get
	r.Get("/cards/today", controllers.GetAllTodayCard)                   // Get all Today's card
	r.Get("/cards/today/one", controllers.GetTodayCard)                  // Get Today card
	r.Get("/cards/next", controllers.GetNextCard)                        // Get Next card
	r.Get("/cards/:deckID/next", controllers.GetNextCardByDeck)          // Get Next card by deck
	r.Get("/cards/:deckID/training", controllers.GetTrainingCardsByDeck) // Get training card by deck

	r.Get("/mcqs/:deckID", controllers.GetMcqsByDeck) // Get MCQs by deckID

	// Post
	r.Post("/cards/response", controllers.PostResponse) // Post a response

	// ADMIN ONLY
	r.Get("/cards", controllers.GetAllCards)                   // Get all cards
	r.Get("/cards/id/:id", controllers.GetCardByID)            // Get card by ID
	r.Get("/cards/deck/:deckID", controllers.GetCardsFromDeck) // Get card by deckID

	r.Post("/cards/new", controllers.CreateNewCard) // Create a new card
	r.Post("/mcqs/new", controllers.CreateMcq)      // Create a mcq

	r.Put("/cards/:id/edit", controllers.UpdateCardByID) // Update a card by ID

	r.Delete("/cards/:id", controllers.DeleteCardById) // Delete a card by ID
}
