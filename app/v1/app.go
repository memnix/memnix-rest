package v1

import (
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"
	_ "github.com/memnix/memnix-rest/docs" // Side effect import
	"github.com/memnix/memnix-rest/pkg/json"
)

// New returns a new Fiber instance.
func New() *fiber.App {
	// Create new app

	app := fiber.New(
		fiber.Config{
			Prefork:     false,
			JSONDecoder: json.NewJSON(&json.SonicJSON{}).Unmarshal,
			JSONEncoder: json.NewJSON(&json.SonicJSON{}).Marshal,
		})

	// Register middlewares
	registerMiddlewares(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// Use swagger middleware
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// Api group
	v1 := app.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	registerRoutes(&v1)

	return app
}

func registerMiddlewares(app *fiber.App) {
	// Use cors middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost, *",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cache-Control",
		AllowCredentials: true,
	}))

	// Provide a minimal config
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
		URL:  "/favicon.ico",
	}))

	app.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))

	app.Use(otelfiber.Middleware(otelfiber.WithNext(
		func(c *fiber.Ctx) bool {
			// Do not trace /metrics endpoint
			return c.Path() == "/metrics" || c.Path() == "/swagger/*" || c.Path() == "/favicon.ico"
		})))

	app.Use(pprof.New())
}
