package http

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/pprof"
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

	app.Use(pprof.New())

	app.Use(loggerMiddleware())

}

func loggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue stack
		chainErr := c.Next()

		if chainErr != nil {
			if err := c.App().ErrorHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError) //nolint:errcheck // TODO: Explain why we ignore the error here
			}
		}

		// Do something with response
		p := influxdb2.NewPointWithMeasurement("fiber").
			AddField("ip", c.IP()).
			AddField("method", c.Method()).
			AddField("path", c.Path()).
			AddField("status", c.Response().StatusCode()).
			SetTime(time.Now())

		misc.LogWriter{}.Write(*p)

		return nil

	}
}
