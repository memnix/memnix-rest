package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/views"
)

type PageController struct{}

func NewPageController() *PageController {
	return &PageController{}
}

func (p *PageController) GetIndex(c echo.Context) error {
	page := views.Page("John")

	return Render(c, http.StatusOK, page)
}

func (p *PageController) PostClicked(c echo.Context) error {
	clicked := views.Clicked()

	return Render(c, http.StatusOK, clicked)
}
