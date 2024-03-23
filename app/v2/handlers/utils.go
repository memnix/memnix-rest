package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/domain"
)

func Render(c echo.Context, _ int, t templ.Component) error {
	c.Response().Writer.Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response())
}

func Redirect(c echo.Context, path string, statusCode int) error {
	c.Response().Header().Set("HX-Redirect", path)
	return c.NoContent(statusCode)
}

func GetNonce(c echo.Context) domain.Nonce {
	return c.Get("nonce").(domain.Nonce)
}
