package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/memnix/memnix-rest/views"
)

type AuthController struct {
	auth auth.IUseCase
}

func NewAuthController(auth auth.IUseCase) AuthController {
	return AuthController{auth: auth}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var loginStruct domain.Login
	err := c.BodyParser(&loginStruct)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	jwtToken, err := a.auth.Login(loginStruct.Password, loginStruct.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(auth.NewLoginTokenVM("", "invalid credentials"))
	}

	return c.Status(fiber.StatusOK).JSON(auth.NewLoginTokenVM(jwtToken, ""))
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	var registerStruct domain.Register
	err := c.BodyParser(&registerStruct)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	newUser, err := a.auth.Register(registerStruct)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(auth.NewRegisterVM("user created", newUser.ToPublicUser()))
}

func (a *AuthController) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(auth.NewLoginTokenVM("", "logged out"))
}

func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	newToken, err := a.auth.RefreshToken(utils.GetUserFromContext(c))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVMFromError(err))
	}

	return c.Status(fiber.StatusOK).JSON(auth.NewLoginTokenVM(newToken, ""))
}
