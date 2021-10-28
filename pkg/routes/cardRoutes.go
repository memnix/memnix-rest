package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerCardRoutes(r fiber.Router) {
	// Get
	r.Get("/cards", controllers.GetAllCards)                   // Get all cards
	r.Get("/cards/id/:id", controllers.GetCardByID)            // Get card by ID
	r.Get("/cards/deck/:deckID", controllers.GetCardsFromDeck) // Get card by deckID

	r.Get("/cards/today", controllers.GetTodayCard) // Get Today card
	// Post
	r.Post("/cards/new", controllers.CreateNewCard) // Create a new deck
}
