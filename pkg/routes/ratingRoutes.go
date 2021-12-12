package routes

import (
	"memnixrest/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func registerRatingRoutes(r fiber.Router) {

	// Get
	r.Get("/ratings", controllers.GetAllRatings) // Get all ratings
	r.Get("/ratings/deck/:deckID", controllers.GetAllRatingsByDeck)
	r.Get("/ratings/deck/:deckID/average", controllers.GetAverageRatingByDeck)
	r.Get("/ratings/deck/:deckID/user/:userID", controllers.GetRatingByDeckAndUser)
	r.Get("/ratings/deck/:deckID/user", controllers.GetRatingByDeckAndUser)

	// Post
	r.Post("/ratings/new", controllers.RateDeck) // Rate a deck
}
