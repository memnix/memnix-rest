package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/views/components"
	"github.com/memnix/memnix-rest/domain"
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
	loginStruct := domain.Login{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	err := validator.New().Struct(loginStruct)
	if err != nil {
		slog.Debug("Auth: ", slog.String("error", err.Error()))
		return Render(c, http.StatusForbidden, components.LoginError("Invalid email or password"))
	}

	slog.Debug("Auth: ", slog.String("email", loginStruct.Email), slog.String("password", loginStruct.Password))

	// Call the use case to authenticate the user
	jwtToken, err := a.useCase.Login(context.Background(), loginStruct.Password, loginStruct.Email)
	if err != nil {
		slog.Debug("Auth: ", slog.String("error", err.Error()))
		return Render(c, http.StatusForbidden, components.LoginError("Invalid email or password"))
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

	setFlashmessages(c, "success", "You are now logged in")

	return Redirect(c, "/", http.StatusAccepted)
}

func (a *AuthController) PostLogout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     SessionTokenCookieKey,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusAccepted)
}

func (a *AuthController) PostRegister(c echo.Context) error {
	registerStruct := domain.Register{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		Username: c.FormValue("username"),
	}

	err := validator.New().Struct(registerStruct)
	if err != nil {
		slog.Debug("Auth: ", slog.String("error", err.Error()))
		return Render(c, http.StatusForbidden, components.UsernameError("Something went wrong, please try again."))
	}

	slog.Debug("Auth: ", slog.String("email", registerStruct.Email),
		slog.String("password", registerStruct.Password),
		slog.String("username", registerStruct.Username))

	// Validate the password
	entropy, err := a.useCase.ValidatePassword(c.Request().Context(), registerStruct.Password)
	if err != nil {
		slog.Info("Auth: ", slog.String("error", err.Error()))
		passwordEntropy := components.PasswordEntropy(entropy, err)
		return Render(c, http.StatusOK, passwordEntropy)
	}

	// Call the use case to authenticate the user
	_, err = a.useCase.Register(c.Request().Context(), registerStruct)
	if err != nil {
		slog.Info("Auth: ", slog.String("error", err.Error()))
		setFlashmessages(c, "error",
			"The email is probably already in use. Please try again with a different email  or login.")
		return Redirect(c, "/register", http.StatusForbidden)
	}

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusAccepted)
}

func (a *AuthController) ValidatePassword(c echo.Context) error {
	password := c.FormValue("password")

	entropy, err := a.useCase.ValidatePassword(c.Request().Context(), password)

	passwordEntropy := components.PasswordEntropy(entropy, err)
	return Render(c, http.StatusOK, passwordEntropy)
}
