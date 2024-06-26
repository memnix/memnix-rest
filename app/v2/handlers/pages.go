package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/views/page"
)

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

func (p *PageController) GetIndex(c echo.Context) error {
	hero := page.Hero("John")

	index := page.HomePage("Memnix", "", false, GetNonce(c), hero)

	return Render(c, http.StatusOK, index)
}

func (p *PageController) PostClicked(c echo.Context) error {
	clicked := page.Clicked()

	return Render(c, http.StatusOK, clicked)
}

func (p *PageController) GetLogin(c echo.Context) error {
	loginContent := page.LoginContent()
	login := page.LoginPage("Login", GetNonce(c), loginContent)

	return Render(c, http.StatusOK, login)
}

func (p *PageController) GetRegister(c echo.Context) error {
	registerContent := page.RegisterContent()
	register := page.RegisterPage("Register", GetNonce(c), registerContent)

	return Render(c, http.StatusOK, register)
}
