package v2

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/memnix/memnix-rest/app/v2/handlers"
	"github.com/memnix/memnix-rest/assets"
	"github.com/memnix/memnix-rest/services"
)

func (i *InstanceSingleton) registerStaticRoutes(e *echo.Echo) {

	g := e.Group("/static", StaticAssetsCacheControlMiddleware)

	g.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".",
		Browse:     false,
		Filesystem: assets.Assets(),
	}))

}

func (i *InstanceSingleton) registerRoutes(e *echo.Echo) {
	serviceContainer := services.DefaultServiceContainer()
	authController := serviceContainer.AuthHandler()
	pageController := handlers.NewPageController()

	// Page routes
	e.GET("/", pageController.GetIndex, StaticPageCacheControlMiddleware)
	e.GET("/login", pageController.GetLogin, StaticPageCacheControlMiddleware)
	e.GET("/register", pageController.GetRegister, StaticPageCacheControlMiddleware)

	// Auth routes
	e.POST("/register", authController.PostRegister)
	e.POST("/logout", authController.PostLogout)
	e.POST("/login", authController.PostLogin)
	e.POST("/register/password", authController.ValidatePassword)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
