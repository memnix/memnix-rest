package v2

import (
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/handlers"
)

func (i *InstanceSingleton) registerStaticRoutes(e *echo.Echo) {
	e.Static("/", "static")
}

func (i *InstanceSingleton) registerRoutes(e *echo.Echo) {
	pageController := handlers.NewPageController()

	e.GET("/", pageController.GetIndex)
	e.POST("/clicked", pageController.PostClicked)
}
