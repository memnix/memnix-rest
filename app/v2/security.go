package v2

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	db "github.com/memnix/memnix-rest/db/sqlc"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/memnix/memnix-rest/services/user"
)

// JwtMiddleware is the controller for the jwt routes.
type JwtMiddleware struct {
	user.IUseCase
}

// NewJwtController creates a new jwt controller.
func NewJwtMiddleware(user user.IUseCase) JwtMiddleware {
	return JwtMiddleware{IUseCase: user}
}

// VerifyPermissions checks if the user has the required permissions.
func (*JwtMiddleware) VerifyPermissions(user db.User, p domain.Permission) bool {
	return domain.Permission(user.Permission.Int32) >= p
}

func (j *JwtMiddleware) AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return j.IsConnectedMiddleware(domain.PermissionUser, next)(c)
	}
}

func (j *JwtMiddleware) AuthorizeVip(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return j.IsConnectedMiddleware(domain.PermissionVip, next)(c)
	}
}

func (j *JwtMiddleware) AuthorizeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return j.IsConnectedMiddleware(domain.PermissionAdmin, next)(c)
	}
}

func (j *JwtMiddleware) IsConnectedMiddleware(p domain.Permission, next echo.HandlerFunc) func(c echo.Context) error {
	return func(c echo.Context) error {
		// get the token from the request header
		tokenHeader := c.Request().Header.Get("Authorization")
		// if the token is empty, the userModel is not connected, and we return an error
		if tokenHeader == "" {
			return c.JSON(http.StatusUnauthorized, errors.New("unauthorized: token missing"))
		}

		userID, err := jwt.GetJwtInstance().GetJwt().GetConnectedUserID(c.Request().Context(), tokenHeader)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, errors.New("unauthorized: invalid token"))
		}

		userModel, err := j.IUseCase.GetByID(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, errors.New("unauthorized: invalid user"))
		}

		// Check permissions
		if !j.VerifyPermissions(userModel, p) {
			return c.JSON(http.StatusUnauthorized, errors.New("unauthorized: invalid user"))
		}

		// Set userModel in locals
		SetUserToContext(c, userModel)
		return next(c)
	}
}

func SetUserToContext(c echo.Context, user db.User) {
	c.Set("user", user)
}

func GetUserFromContext(c echo.Context) (db.User, error) {
	if c.Get("user") == nil {
		return db.User{}, errors.New("user is not initialized")
	}
	return c.Get("user").(db.User), nil
}
