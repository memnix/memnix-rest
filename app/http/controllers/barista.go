package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal/barista"
)

type BaristaController struct {
	barista.IUseCase
}

func NewBaristaController(useCase barista.IUseCase) BaristaController {
	return BaristaController{IUseCase: useCase}
}

func (b *BaristaController) GetName(c *fiber.Ctx) error {
	return c.SendString(b.IUseCase.GetName())
}
