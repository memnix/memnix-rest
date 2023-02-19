package controllers

import (
	"github.com/edgedb/edgedb-go"
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
	// Convert uuid to edgedb.UUID
	uuidEdge, err := edgedb.ParseUUID(uuid)
	if err != nil {
		return err
	}

	return c.SendString(u.IUseCase.GetName(uuidEdge))
}
