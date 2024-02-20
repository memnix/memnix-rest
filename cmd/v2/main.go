package main

import (
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/views"
)

func main() {
	e := echo.New()
	component := views.Page("John")

	e.Static("/", "static")

	e.GET("/", func(c echo.Context) error {
		err := component.Render(c.Request().Context(), c.Response())
		if err != nil {
			return err
		}
		return nil
	})

	e.Logger.Fatal(e.Start(":3000"))
}
