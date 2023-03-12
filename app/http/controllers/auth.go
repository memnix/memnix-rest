package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/auth"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/memnix/memnix-rest/views"
	"github.com/pkg/errors"
)

// AuthController is the controller for the auth routes
type AuthController struct {
	auth auth.IUseCase
}

// NewAuthController creates a new auth controller
func NewAuthController(auth auth.IUseCase) AuthController {
	return AuthController{auth: auth}
}

// Login is the controller for the login route
//
//	@Summary		Login
//	@Description	Login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body		domain.Login	true	"Login"
//	@Success		200		{object}	views.LoginTokenVM
//	@Failure		401		{object}	views.HTTPResponseVM
//	@Failure		500		{object}	views.HTTPResponseVM
//	@Router			/v2/security/login [post]
func (a *AuthController) Login(c *fiber.Ctx) error {
	var loginStruct domain.Login
	err := c.BodyParser(&loginStruct)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	jwtToken, err := a.auth.Login(loginStruct.Password, loginStruct.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewLoginTokenVM("", "invalid credentials"))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewLoginTokenVM(jwtToken, ""))
}

// Register is the controller for the register route
//
//	@Summary		Register
//	@Description	Register
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			register	body		domain.Register	true	"Register"
//	@Success		201			{object}	views.RegisterVM
//	@Failure		400			{object}	views.HTTPResponseVM
//	@Failure		500			{object}	views.HTTPResponseVM
//	@Router			/v2/security/register [post]
func (a *AuthController) Register(c *fiber.Ctx) error {
	var registerStruct domain.Register
	err := c.BodyParser(&registerStruct)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(err))
	}

	newUser, err := a.auth.Register(registerStruct)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(views.NewHTTPResponseVMFromError(errors.New("error creating user")))
	}

	return c.Status(fiber.StatusCreated).JSON(views.NewRegisterVM("user created", newUser.ToPublicUser()))
}

// Logout is the controller for the logout route
//
//	@Summary		Logout
//	@Description	Logout
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	views.LoginTokenVM
//	@Failure		500	{object}	views.HTTPResponseVM
//	@Router			/v2/security/logout [post]
func (a *AuthController) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(views.NewLoginTokenVM("", "logged out"))
}

// RefreshToken is the controller for the refresh token route
//
//	@Summary		Refresh token
//	@Description	Refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	views.LoginTokenVM
//	@Failure		401	{object}	views.HTTPResponseVM
//	@Failure		500	{object}	views.HTTPResponseVM
//	@Router			/v2/security/refresh [post]
func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	newToken, err := a.auth.RefreshToken(*utils.GetUserFromContext(c))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVMFromError(err))
	}

	return c.Status(fiber.StatusOK).JSON(views.NewLoginTokenVM(newToken, ""))
}

// SearchKey function to search the key
//
//	@Description	Search the key
//	@Summary		searches the key
//	@Tags			security
//	@Produce		json
//	@Success		200	{object}	views.HTTPResponseVM
//	@Failure		401	{object}	views.HTTPResponseVM
//	@Router			/v2/security/key [get]
//	@Security		Beaver
func (a *AuthController) SearchKey(c *fiber.Ctx) error {
	// Search the key
	key, err := infrastructures.GetSearchKey()
	if err != nil {
		return err
	}

	// Return the keys
	return c.Status(fiber.StatusOK).JSON(views.NewLoginTokenVM(key.Key, ""))
}
