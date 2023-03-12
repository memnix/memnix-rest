package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/memnix/memnix-rest/views"
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

func (u *UserController) GetMe(c *fiber.Ctx) error {
	userCtx := utils.GetUserFromContext(c)
	if userCtx == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVM("User not found", nil))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewHTTPResponseVM("User found", *userCtx))
}
