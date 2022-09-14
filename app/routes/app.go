package routes

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/memnix/memnixrest/docs" // Side effect import
	"github.com/memnix/memnixrest/pkg/models"

	"time"
)

type routeStruct struct {
	Method     string
	Handler    func(c *fiber.Ctx) error
	Permission models.Permission
}

var routesMap map[string]routeStruct

func New() *fiber.App {
	// Create new app
	app := fiber.New(
		fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost, *",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 2
	}))

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true" || c.Path() == "/v1/user" || c.Path() == "/v1/login" || c.Path() == "/v1/register" || c.Path() == "/v1/logout"
		},
		Expiration:   2 * time.Minute,
		CacheControl: true,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// Api group
	v1 := app.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "This is not a valid route") // Custom error
	})

	v1.Use(IsConnectedMiddleware())

	// Register routes
	routesMap = make(map[string]routeStruct)

	registerAuthRoutes() // /v1/
	registerUserRoutes() // /v1/users/
	registerDeckRoutes() // /v1/decks/
	registerCardRoutes() // /v1/cards/

	for route, routeData := range routesMap {
		v1.Add(routeData.Method, route, routeData.Handler)
	}

	return app
}
