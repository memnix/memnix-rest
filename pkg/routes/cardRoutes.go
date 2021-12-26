package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerCardRoutes(r fiber.Router) {
	// Get
	r.Get("/cards/today", controllers.GetTodayCard)             // Get Today card
	r.Get("/cards/next", controllers.GetNextCard)               // Get Next card
	r.Get("/cards/:deckID/next", controllers.GetNextCardByDeck) // Get Next card by deck

	// Post
	r.Post("/cards/response", controllers.PostResponse) // Post a response

	// ADMIN ONLY
	r.Get("/cards", controllers.GetAllCards)                   // Get all cards
	r.Get("/cards/id/:id", controllers.GetCardByID)            // Get card by ID
	r.Get("/cards/deck/:deckID", controllers.GetCardsFromDeck) // Get card by deckID

	r.Post("/cards/new", controllers.CreateNewCard)                   // Create a new card
	r.Post("/cards/deck/:deckID/bulk", controllers.CreateNewCardBulk) // Create a list of cards

}
