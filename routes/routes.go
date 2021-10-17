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
	// Get
	v1.Get("/users", handlers.GetAllUsers)                          // Get all users
	v1.Get("/users/id/:id", handlers.GetUserByID)                   // Get user by ID
	v1.Get("users/discord/:discordID", handlers.GetUserByDiscordID) // Get user by discordID
	// Post
	v1.Post("/users/new", handlers.CreateNewUser) // Create a new user
	// Put
	v1.Put("/users/id/:id", handlers.UpdateUserByID) // Update an user using his ID

	// Decks
	// Get
	v1.Get("/decks", handlers.GetAllDecks)        // Get all decks
	v1.Get("/decks/id/:id", handlers.GetDeckByID) // Get deck by ID
	// Post
	v1.Post("/decks/new", handlers.CreateNewDeck)                        // Create a new deck
	v1.Post("/decks/:deckID/user/:userID/subscribe", handlers.SubToDeck) // Subscribe to a deck

	// Cards
	// Get
	v1.Get("/cards", handlers.GetAllCards)                   // Get all cards
	v1.Get("/cards/id/:id", handlers.GetCardByID)            // Get card by ID
	v1.Get("/cards/deck/:deckID", handlers.GetCardsFromDeck) // Get card by deckID
	// Post
	v1.Post("/cards/new", handlers.CreateNewCard)    // Create a new deck
	debug.Get("/cards", handlers.GetRandomDebugCard) // DEBUG: Get random card

	// Mem
	// Get
	v1.Get("/mems/user/:userID/deck/:deckID/next", handlers.GetNextMem)       // Get next mem by userID & deckID
	v1.Get("/mems/user/:userID/deck/:deckID/today", handlers.GetTodayNextMem) // Get today's next mem by userID & deckID
	v1.Get("/mems/id/:id", handlers.GetMemByID)                               // Get mem by ID
	v1.Get("/mems/user/:userID/card/:cardID", handlers.GetMemByCardAndUser)   // Get mem by userID & cardID
	// Post
	v1.Post("/mems/new", handlers.CreateNewMem) // Create a new mem
	// Put
	v1.Put("/mem/id/:id", handlers.UpdateMemByID) // Update mem by ID

	// Revision
	// Get
	v1.Get("/revisions", handlers.GetAllRevisions)                  // Get all revisions
	v1.Get("/revisions/id/:id", handlers.GetRevisionByID)           // Get revision by ID
	v1.Get("/revisions/user/:userID", handlers.GetRevisionByUserID) // Get revision by userID
	// Post
	v1.Post("/revisions/new", handlers.CreateNewRevision) // Create a new revision

	// Access
	// Get
	v1.Get("/accesses", handlers.GetAllAccesses)                                       // Get all accesses
	v1.Get("/accesses/id/:id", handlers.GetAccessByID)                                 // Get access by ID
	v1.Get("/accesses/user/:userID/deck/:deckID", handlers.GetAccessByUserIDAndDeckID) // Get access by userID & deckID
	v1.Get("/accesses/user/:userID", handlers.GetAccessesByUserID)                     // Get accesses by userID
	// Post
	v1.Post("/accesses/new", handlers.CreateNewAccess) // Create a new access
	// Put
	v1.Put("/accesses/id/:id", handlers.UpdateAccessByID) // Update an access using his ID

	// Answer
	// Get
	v1.Get("/answers", handlers.GetAllAnswers)                   // Get all answers
	v1.Get("/answers/id/:id", handlers.GetAnswerByID)            // Get answer by ID
	v1.Get("/answers/card/:cardID", handlers.GetAnswersByCardID) // Get answer by CardID
	// Post
	v1.Post("/answers/new", handlers.CreateNewAnswer) // Create a new answer

	// History
	// TODO

	return app
}
