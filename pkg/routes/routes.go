package routes

import (
	"memnixrest/app/controllers"

	_ "memnixrest/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	// Create new app
	app := fiber.New()

	app.Get("/swagger/*", swagger.Handler) // default

	// Api group
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	// Auth
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/user", controllers.User)
	api.Post("/logout", controllers.Logout)

	// v1 group "/api/v1"
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		return c.Next()
	})

	v1.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	// Register routes
	registerUserRoutes(v1) // /v1/users/
	registerDeckRoutes(v1) // /v1/decks/
	registerCardRoutes(v1) // /v1/cards/

	// Mem
	// Get
	v1.Get("/mems/id/:id", controllers.GetMemByID)                             // Get mem by ID
	v1.Get("/mems/user/:userID/card/:cardID", controllers.GetMemByCardAndUser) // Get mem by userID & cardID
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

	return app
}
