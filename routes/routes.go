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
	v1.Get("/users", handlers.GetAllUsers)                         // Get all users
	v1.Get("/user/id/:id", handlers.GetUserByID)                   // Get user by ID
	v1.Get("user/discord/:discordID", handlers.GetUserByDiscordID) // Get user by discordID
	v1.Post("/user/new", handlers.CreateNewUser)                   // Create a new user
	v1.Put("/user/id/:id", handlers.UpdateUserByID)                // Update an user using his ID

	// Decks
	v1.Get("/decks", handlers.GetAllDecks)       // Get all decks
	v1.Get("/deck/id/:id", handlers.GetDeckByID) // Get deck by ID
	v1.Post("/deck/new", handlers.CreateNewDeck) // Create a new deck

	// Cards
	v1.Get("/cards", handlers.GetAllCards)                  // Get all cards
	v1.Get("/card/id/:id", handlers.GetCardByID)            // Get card by ID
	v1.Get("/card/deck/:deckID", handlers.GetCardsFromDeck) // Get card by deckID
	v1.Post("/card/new", handlers.CreateNewCard)            // Create a new deck
	debug.Get("/card", handlers.GetRandomDebugCard)         // DEBUG: Get random card

	// Mem
	debug.Get("/user/:userID/deck/:deckID/next", handlers.GetNextCard)        // Get next mem by userID & deckID
	debug.Get("/user/:userID/deck/:deckID/today", handlers.GetTodayNextCard)  // Get today's next mem by userID & deckID
	debug.Get("/mem/id/:id", handlers.GetMemByID)                             // Get mem by ID
	debug.Get("/mem/user/:userID/card/:cardID", handlers.GetMemByCardAndUser) // Get mem by userID & cardID
	debug.Put("/mem/id/:id", handlers.UpdateMemByID)                          // Update mem by ID
	debug.Post("/mem/new", handlers.CreateNewMem)                             // Create a new mem
	debug.Post("/deck/:deckID/sub", handlers.SubToDeck)                       // Subscribe to a deck

	// Revision
	v1.Get("/revisions", handlers.GetAllRevisions)                 // Get all revisions
	v1.Get("/revision/id/:id", handlers.GetRevisionByID)           // Get revision by ID
	v1.Get("/revision/user/:userID", handlers.GetRevisionByUserID) // Get revision by userID
	v1.Get("/revision/card/cardID", handlers.GetRevisionByCardID)  // Get revision by cardID
	v1.Post("/revision/new", handlers.CreateNewRevision)           // Create a new revision

	// Access
	v1.Get("/accesses", handlers.GetAllAccesses)                                   // Get all users
	v1.Get("/access/id/:id", handlers.GetAccessByID)                               // Get user by ID
	v1.Get("/access/user/:userID/deck/:deckID", handlers.GetAccessByUserAndDeckID) // Get user by ID
	v1.Get("accesses/user/:userID", handlers.GetAccessesByUserID)                  // Get accesses by userID
	v1.Post("/access/new", handlers.CreateNewAccess)                               // Create a new user
	v1.Put("/access/id/:id", handlers.UpdateAccessByID)                            // Update an user using his ID

	// History
	// TODO

	return app

}
