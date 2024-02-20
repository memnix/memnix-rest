package controllers

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/pkg/json"
	"github.com/pkg/errors"
)

// GetUserFromContext gets the user from the context.
func GetUserFromContext(ctx *fiber.Ctx) (domain.User, error) {
	if ctx.Locals("user") == nil {
		log.WithContext(ctx.UserContext()).Error("User not found in context")
		return domain.User{}, errors.New("User is not initialized")
	}
	return ctx.Locals("user").(domain.User), nil
}

// SetUserToContext sets the user to the context.
func SetUserToContext(ctx *fiber.Ctx, user domain.User) {
	ctx.Locals("user", user)
}

type JSONHelperSingleton struct {
	jsonHelper json.Helper
}

var (
	instance *JSONHelperSingleton //nolint:gochecknoglobals //Singleton
	once     sync.Once            //nolint:gochecknoglobals //Singleton
)

func GetJSONHelperInstance() *JSONHelperSingleton {
	once.Do(func() {
		instance = &JSONHelperSingleton{
			jsonHelper: json.NewJSON(&json.SonicJSON{}),
		}
	})
	return instance
}

func (j *JSONHelperSingleton) GetJSONHelper() json.Helper {
	return j.jsonHelper
}

func (j *JSONHelperSingleton) SetJSONHelper(helper json.Helper) {
	j.jsonHelper = json.NewJSON(helper)
}
