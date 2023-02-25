package http

import (
	"time"

	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"
	"github.com/memnix/memnix-rest/app/misc"
	"github.com/memnix/memnix-rest/config"
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

	// Use monitor middleware
	app.Get("/metrics", monitor.New(monitor.Config{
		Title:   "Kafejo API",
		Refresh: time.Second * 5,
	}))

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

	// Use compress middleware
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	app.Use(pprof.New())

	cfg := fibernewrelic.Config{
		Application: infrastructures.GetRelicApp(),
	}

	app.Use(fibernewrelic.New(cfg))

	// User logging middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] - [${ip}]:${port} - ${latency} ${method} ${path} - ${status}\n",
		TimeFormat: "Jan 02 | 15:04:05",
		Output:     misc.LogWriter{},
	}))

	app.Use(pprof.New())
}
