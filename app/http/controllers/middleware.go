package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal/user"
	"github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/memnix/memnix-rest/views"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// JwtController is the controller for the jwt routes
type JwtController struct {
	user.IUseCase
}

// NewJwtController creates a new jwt controller
func NewJwtController(user user.IUseCase) JwtController {
	return JwtController{IUseCase: user}
}

// VerifyPermissions checks if the user has the required permissions
func (*JwtController) VerifyPermissions(user domain.User, p domain.Permission) bool {
	return user.HasPermission(p)
}

// IsConnectedMiddleware checks if the user is connected and has the required permissions
// the permissions are defined in the route definition
// returns an error if the user is not connected or has not the required permissions
//
// if the user is connected and has the required permissions, it sets the user in the locals
// and calls the next middleware
func (j *JwtController) IsConnectedMiddleware(p domain.Permission) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// if the route is public, we don't need to check if the userModel is connected
		if p == domain.PermissionNone {
			return c.Next()
		}

		_, span := infrastructures.GetFiberTracer().Start(c.UserContext(), "IsConnectedMiddleware")
		defer span.End()

		// get the token from the request header
		tokenHeader := c.Get("Authorization")
		// if the token is empty, the userModel is not connected, and we return an error
		if tokenHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVMFromError(errors.New("not authorized")))
		}

		// get the userModel from the token
		// if the token is invalid, we return an error
		userID, err := jwt.GetConnectedUserID(c.UserContext(), tokenHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVMFromError(errors.New("not connected")))
		}

		// get the userModel from the database
		userModel, err := j.IUseCase.GetByID(c.UserContext(), userID)
		if err != nil {
			otelzap.Ctx(c.UserContext()).Error("error getting user / not connected", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVMFromError(errors.New("not connected")))
		}

		// if the userModel has the required permissions, we set the userModel in the locals and call the next middleware
		if j.VerifyPermissions(userModel, p) {
			utils.SetUserToContext(c, userModel) // Set userModel in locals
			span.End()
			return c.Next()
		}

		// if the userModel does not have the required permissions, we return an error
		otelzap.Ctx(c.UserContext()).Warn("not authorized", zap.Error(errors.New("not authorized")))
		return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVMFromError(errors.New("not authorized")))
	}
}
