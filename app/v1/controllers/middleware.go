package controllers

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/memnix/memnix-rest/app/v1/views"
	"github.com/memnix/memnix-rest/cmd/v1/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/services/user"
)

// JwtController is the controller for the jwt routes.
type JwtController struct {
	user.IUseCase
}

// NewJwtController creates a new jwt controller.
func NewJwtController(user user.IUseCase) JwtController {
	return JwtController{IUseCase: user}
}

// VerifyPermissions checks if the user has the required permissions.
func (*JwtController) VerifyPermissions(user domain.User, p domain.Permission) bool {
	return user.HasPermission(p)
}

// IsConnectedMiddleware checks if the user is connected and has the required permissions
// the permissions are defined in the route definition
// returns an error if the user is not connected or has not the required permissions
//
// if the user is connected and has the required permissions, it sets the user in the locals
// and calls the next middleware.
func (j *JwtController) IsConnectedMiddleware(p domain.Permission) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// check if the permission is valid
		if !p.IsValid() {
			return c.Status(fiber.StatusInternalServerError).JSON(
				views.NewHTTPResponseVMFromError(errors.New("invalid permission")))
		}

		// if the route is public, we don't need to check if the userModel is connected
		if p == domain.PermissionNone {
			return c.Next()
		}

		_, span := infrastructures.GetTracerInstance().Tracer().Start(c.UserContext(), "IsConnectedMiddleware")
		defer span.End()

		// get the token from the request header
		tokenHeader := c.Get("Authorization")
		// if the token is empty, the userModel is not connected, and we return an error
		if tokenHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				views.NewHTTPResponseVMFromError(errors.New("unauthorized: token missing")))
		}

		// get the userModel from the token
		// if the token is invalid, we return an error
		userID, err := config.GetJwtInstance().GetJwt().GetConnectedUserID(c.UserContext(), tokenHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				views.NewHTTPResponseVMFromError(errors.New("unauthorized: invalid token")))
		}

		// get the userModel from the database
		userModel, err := j.IUseCase.GetByID(c.UserContext(), userID)
		if err != nil {
			log.WithContext(c.UserContext()).Error("error getting user / not connected", slog.Any("error", err))
			return c.Status(fiber.StatusUnauthorized).JSON(
				views.NewHTTPResponseVMFromError(errors.New("unauthorized: invalid user")))
		}

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{
				ID:       strconv.Itoa(int(userModel.ID)),
				Username: userModel.Username,
				Email:    userModel.Email,
			})
		})

		// Check permissions
		if !j.VerifyPermissions(userModel, p) {
			log.WithContext(c.UserContext()).Debug("Not authorized")
			return c.Status(fiber.StatusUnauthorized).JSON(views.NewHTTPResponseVMFromError(errors.New("unauthorized: insufficient permissions")))
		}

		// Set userModel in locals
		SetUserToContext(c, userModel)
		span.End()
		return c.Next()
	}
}
