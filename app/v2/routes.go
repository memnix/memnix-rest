package v2

import (
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/handlers"
	"github.com/memnix/memnix-rest/services"
)

func (i *InstanceSingleton) registerStaticRoutes(e *echo.Echo) {
	e.Static("/", "assets/static")
	e.Static("/img", "assets/img")
}

func (i *InstanceSingleton) registerRoutes(e *echo.Echo) {
	serviceContainer := services.DefaultServiceContainer()
	authController := serviceContainer.AuthHandler()
	pageController := handlers.NewPageController()

	e.GET("/", pageController.GetIndex)
	e.GET("/login", pageController.GetLogin)
	e.POST("/login", authController.PostLogin)
	e.POST("/clicked", pageController.PostClicked)
}
