package http

import (
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/memnix/memnix-rest/app/misc"
	"github.com/memnix/memnix-rest/config"
	_ "github.com/memnix/memnix-rest/docs" // Side effect import
	"github.com/memnix/memnix-rest/infrastructures"
)

// New returns a new Fiber instance
func New() *fiber.App {
	// Create new app

	app := fiber.New(
		fiber.Config{
			Prefork:     false,
			JSONDecoder: config.JSONHelper.Unmarshal,
			JSONEncoder: config.JSONHelper.Marshal,
		})

	// Register middlewares
	registerMiddlewares(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// Use swagger middleware
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// Api group
	v2 := app.Group("/v2")

	v2.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	registerRoutes(&v2) // /v2

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

	app.Use(cache.New(cache.Config{
		Expiration:   5 * time.Second,
		CacheControl: true,
		Next: func(c *fiber.Ctx) bool {
			// Do not cache /metrics endpoint
			return c.Path() == "/metrics"
		},
	}))

	cfg := fibernewrelic.Config{
		Application: infrastructures.GetRelicApp(),
	}

	app.Use(fibernewrelic.New(cfg))

	prometheus := fiberprometheus.New("memnix")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// User logging middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] - [${ip}]:${port} - ${latency} ${method} ${path} - ${status}\n",
		TimeFormat: "Jan 02 | 15:04:05",
		Output:     misc.LogWriter{},
	}))
}
