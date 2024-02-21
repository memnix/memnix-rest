package v2

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/views"
)

func registerStaticRoutes(e *echo.Echo) {
	e.Static("/", "static")
}

func registerRoutes(e *echo.Echo) {
	component := views.Page("John")

	clickedComponent := views.Clicked()

	e.GET("/", func(c echo.Context) error {
		err := component.Render(c.Request().Context(), c.Response())
		if err != nil {
			return err
		}
		return nil
	})

	e.POST("/clicked", func(c echo.Context) error {
		return Render(c, http.StatusOK, clickedComponent)
	})
}

func Render(c echo.Context, _ int, t templ.Component) error {
	c.Response().Writer.Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response())
}
