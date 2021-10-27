package routes

import (
	"memnixrest/app/controllers"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	_ "memnixrest/docs"
)

func New() *fiber.App {
	// Create new app
	app := fiber.New()

	app.Get("/swagger/*", swagger.Handler) // default

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
	v1.Get("/users", controllers.GetAllUsers)                          // Get all users
	v1.Get("/users/id/:id", controllers.GetUserByID)                   // Get user by ID
	v1.Get("users/discord/:discordID", controllers.GetUserByDiscordID) // Get user by discordID
	// Post
	v1.Post("/users/new", controllers.CreateNewUser) // Create a new user
	// Put
	v1.Put("/users/id/:id", controllers.UpdateUserByID) // Update an user using his ID

	// Decks
	// Get
	v1.Get("/decks", controllers.GetAllDecks)                 // Get all decks
	v1.Get("/decks/public", controllers.GetAllPublicDecks)    // Get all public decks
	v1.Get("/decks/user/:userID", controllers.GetAllSubDecks) // Get all decks the user is sub to
	v1.Get("/decks/id/:id", controllers.GetDeckByID)          // Get deck by ID
	// Post
	v1.Post("/decks/new", controllers.CreateNewDeck)                            // Create a new deck
	v1.Post("/decks/:deckID/user/:userID/subscribe", controllers.SubToDeck)     // Subscribe to a deck
	v1.Post("/decks/:deckID/user/:userID/unsubscribe", controllers.UnSubToDeck) // Unsubscribe to a deck

	// Cards
	// Get
	v1.Get("/cards", controllers.GetAllCards)                   // Get all cards
	v1.Get("/cards/id/:id", controllers.GetCardByID)            // Get card by ID
	v1.Get("/cards/deck/:deckID", controllers.GetCardsFromDeck) // Get card by deckID
	// Post
	v1.Post("/cards/new", controllers.CreateNewCard)    // Create a new deck
	debug.Get("/cards", controllers.GetRandomDebugCard) // DEBUG: Get random card

	// Mem
	// Get
	v1.Get("/mems/user/:userID/deck/:deckID/next", controllers.GetNextMem)       // Get next mem by userID & deckID
	v1.Get("/mems/user/:userID/deck/:deckID/today", controllers.GetTodayNextMem) // Get today's next mem by userID & deckID
	v1.Get("/mems/id/:id", controllers.GetMemByID)                               // Get mem by ID
	v1.Get("/mems/user/:userID/card/:cardID", controllers.GetMemByCardAndUser)   // Get mem by userID & cardID
	// Post
	v1.Post("/mems/new", controllers.CreateNewMem) // Create a new mem
	// Put
	v1.Put("/mem/id/:id", controllers.UpdateMemByID) // Update mem by ID


	// Access
	// Get
	v1.Get("/accesses", controllers.GetAllAccesses)                                       // Get all accesses
	v1.Get("/accesses/id/:id", controllers.GetAccessByID)                                 // Get access by ID
	v1.Get("/accesses/user/:userID/deck/:deckID", controllers.GetAccessByUserIDAndDeckID) // Get access by userID & deckID
	v1.Get("/accesses/user/:userID", controllers.GetAccessesByUserID)                     // Get accesses by userID
	// Post
	v1.Post("/accesses/new", controllers.CreateNewAccess) // Create a new access
	// Put
	v1.Put("/accesses/id/:id", controllers.UpdateAccessByID) // Update an access using his ID

	// Answer
	// Get
	v1.Get("/answers", controllers.GetAllAnswers)                   // Get all answers
	v1.Get("/answers/id/:id", controllers.GetAnswerByID)            // Get answer by ID
	v1.Get("/answers/card/:cardID", controllers.GetAnswersByCardID) // Get answer by CardID
	// Post
	v1.Post("/answers/new", controllers.CreateNewAnswer) // Create a new answer

	// History
	// TODO

	return app
}
