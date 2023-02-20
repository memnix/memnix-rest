package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/internal/auth"
)

type AuthController struct {
	auth auth.IUseCase
}

func NewAuthController(auth auth.IUseCase) AuthController {
	return AuthController{auth: auth}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	return c.SendString("Register")
}

func (a *AuthController) Logout(c *fiber.Ctx) error {
	return c.SendString("Logout")
}

func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	return c.SendString("RefreshToken")
}
