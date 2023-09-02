package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/domain"
	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// GetUserFromContext gets the user from the context
func GetUserFromContext(ctx *fiber.Ctx) (domain.User, error) {
	if ctx.Locals("user") == nil {
		otelzap.Ctx(ctx.UserContext()).Error("User not found in context")
		return domain.User{}, errors.New("User is not initialized")
	}
	return ctx.Locals("user").(domain.User), nil
}

// SetUserToContext sets the user to the context
func SetUserToContext(ctx *fiber.Ctx, user domain.User) {
	ctx.Locals("user", user)
}
