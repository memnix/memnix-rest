package v2

import (
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/handlers"
)

func registerStaticRoutes(e *echo.Echo) {
	e.Static("/", "static")
}

func registerRoutes(e *echo.Echo) {
	pageController := handlers.NewPageController()

	e.GET("/", pageController.GetIndex)
	e.POST("/clicked", pageController.PostClicked)
}
