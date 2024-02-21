package v2

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func registerStaticRoutes(e *echo.Echo) {
	e.Static("/", "static")
}

func registerRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World 👋!")
	})

	e.POST("/clicked", func(c echo.Context) error {
		return c.String(http.StatusOK, "Clicked")
	})
}
