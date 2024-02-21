package v2

import (
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	instance *echo.Echo //nolint:gochecknoglobals //Singleton
	once     sync.Once  //nolint:gochecknoglobals //Singleton

)

// New returns a new Echo instance.
func New() *echo.Echo {
	once.Do(func() {
		instance = echo.New()
		registerMiddlewares(instance)

		registerStaticRoutes(instance)

		registerRoutes(instance)
	})
	return instance
}

func registerMiddlewares(_ *echo.Echo) {
}
