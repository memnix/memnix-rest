package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerRatingRoutes(r fiber.Router) {

	// Get
	r.Get("/ratings", controllers.GetAllRatings)                                    // Get all ratings
	r.Get("/ratings/deck/:deckID", controllers.GetAllRatingsByDeck)                 // Get all ratings by deck
	r.Get("/ratings/deck/:deckID/average", controllers.GetAverageRatingByDeck)      // Get average rating by deck
	r.Get("/ratings/deck/:deckID/user/:userID", controllers.GetRatingByDeckAndUser) // Get rating by deck and user
	r.Get("/ratings/deck/:deckID/user", controllers.GetRatingsByDeck)               // Get rating by deck

	// Post
	r.Post("/ratings/new", controllers.RateDeck) // Rate a deck
}
