package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/services/auth"
)

type AuthController struct {
	useCase auth.IUseCase
}

const (
	// SessionTokenCookieKey is the key for the session token cookie.
	SessionTokenCookieKey = "session_token"
	ExpiresDuration       = 24 * time.Hour
)

func NewAuthController(auth auth.IUseCase) AuthController {
	return AuthController{useCase: auth}
}

func (a *AuthController) PostLogin(c echo.Context) error {
	// Get the username and password from the request
	email := c.FormValue("email")
	password := c.FormValue("password")

	slog.Info("Auth: ", slog.String("email", email), slog.String("password", password))

	if a.useCase == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	// Call the use case to authenticate the user
	jwtToken, err := a.useCase.Login(context.Background(), password, email)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "invalid credentials"})
	}

	cookie := &http.Cookie{
		Name:     SessionTokenCookieKey,
		Value:    jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(ExpiresDuration),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}
