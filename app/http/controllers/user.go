package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal/user"
)

type UserController struct {
	user.IUseCase
}

func NewUserController(useCase user.IUseCase) UserController {
	return UserController{IUseCase: useCase}
}

func (u *UserController) GetName(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	return c.SendString(u.IUseCase.GetName(uuid))
}
