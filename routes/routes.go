package routes

import (
	"memnixrest/handlers"

	"github.com/gofiber/fiber/v2"
)

// New
func New() *fiber.App {
	// Create new app
	app := fiber.New()

	// Api group
	api := app.Group("/api")

	// v1 group "/api/v1"
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		return c.Next()
	})

	// debug group "/api/debug"
	debug := api.Group("/debug", func(c *fiber.Ctx) error {
		return c.Next()
	})

	api.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	v1.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	// Users
	v1.Get("/users", handlers.GetAllUsers)       // Get all users
	v1.Get("/user/id/:id", handlers.GetUserByID) // Get user by id
	v1.Get("user/discordid/:discordID", handlers.GetUserByDiscordID)
	v1.Post("/user/new", handlers.CreateNewUser)    // Create a new user
	v1.Put("/user/id/:id", handlers.UpdateUserByID) // Update an user using his id

	// Decks
	v1.Get("/decks", handlers.GetAllDecks)
	v1.Get("/deck/id/:id", handlers.GetDeckByID)
	v1.Post("/deck/new", handlers.CreateNewDeck)

	// Cards
	v1.Get("/cards", handlers.GetAllCards)
	v1.Get("/card/id/:id", handlers.GetCardByID)
	v1.Get("/card/deck/:deckID", handlers.GetCardsFromDeck)
	v1.Post("/card/new", handlers.CreateNewCard)

	debug.Get("/card", handlers.GetRandomDebugCard)

	// Mem
	debug.Get("/user/:userID/deck/:deckID/next", handlers.GetNextCard)
	debug.Get("/user/:userID/deck/:deckID/today", handlers.GetTodayNextCard)
	debug.Post("/mem/new", handlers.CreateNewMem)
	debug.Put("/mem/id/:id", handlers.UpdateMemByID)
	debug.Get("/mem/id/:id", handlers.GetMemByID)
	debug.Get("/mem/user/:userID/card/:cardID", handlers.GetMemByCardAndUser)
	debug.Post("/deck/:deckID/sub", handlers.SubToDeck)

	// Revision
	v1.Get("/revisions", handlers.GetAllRevisions)
	v1.Get("/revision/id/:id", handlers.GetRevisionByID)
	v1.Get("/revision/userid/:userID", handlers.GetRevisionByUserID)
	v1.Get("/revision/cardid/cardID", handlers.GetRevisionByCardID)
	v1.Post("/revision/new", handlers.CreateNewRevision)

	// Access
	// TODO

	// History
	// TODO

	return app

}
