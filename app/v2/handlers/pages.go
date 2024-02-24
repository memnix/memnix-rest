package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/views/page"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

func (p *PageController) GetIndex(c echo.Context) error {
	if lang, ok := c.Get("lang").(string); ok {
		slog.Debug("🌐 ", slog.String("lang", lang))
	}

	localizer, ok := c.Get("localizer").(*i18n.Localizer)
	if !ok {
		slog.WarnContext(c.Request().Context(), "Failed to get localizer")
		return Render(c, http.StatusInternalServerError, nil)
	}

	hero := page.Hero("John", localizer)

	index := page.HomePage("Memnix", "", false, false, nil, nil, hero)

	return Render(c, http.StatusOK, index)
}

func (p *PageController) PostClicked(c echo.Context) error {
	clicked := page.Clicked()

	return Render(c, http.StatusOK, clicked)
}

func (p *PageController) GetLogin(c echo.Context) error {
	errorMessages := getFlashmessages(c, "error")
	successMessages := getFlashmessages(c, "success")

	slog.Debug("Error messages: ", slog.Any("error", errorMessages))
	slog.Debug("Success messages: ", slog.Any("success", successMessages))
	loginContent := page.LoginContent()
	login := page.LoginPage("Login", false, errorMessages,
		successMessages, loginContent)

	return Render(c, http.StatusOK, login)
}

func (p *PageController) GetRegister(c echo.Context) error {
	registerContent := page.RegisterContent()
	register := page.RegisterPage("Register", false, nil, nil, registerContent)

	return Render(c, http.StatusOK, register)
}
