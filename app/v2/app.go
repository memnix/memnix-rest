package v2

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/memnix/memnix-rest/cmd/v2/config"
	"github.com/memnix/memnix-rest/domain"
)

var (
	instance *InstanceSingleton //nolint:gochecknoglobals //Singleton
	once     sync.Once          //nolint:gochecknoglobals //Singleton
)

type InstanceSingleton struct {
	echoInstance *echo.Echo
	config       config.ServerConfig
}

// New returns a new Echo instance.
func GetEchoInstance() *echo.Echo {
	return instance.echoInstance
}

func GetEchoSingleton() *InstanceSingleton {
	once.Do(func() {
		instance = &InstanceSingleton{}
		instance.echoInstance = echo.New()
		instance.registerMiddlewares(instance.echoInstance)

		instance.registerStaticRoutes(instance.echoInstance)

		instance.registerRoutes(instance.echoInstance)
	})
	return instance
}

func CreateEchoInstance(config config.ServerConfig) *InstanceSingleton {
	return GetEchoSingleton().WithConfig(config)
}

func (i *InstanceSingleton) Start() error {
	if err := i.echoInstance.Start(":" + i.config.Port); err != nil {
		return err
	}

	return nil
}

func (i *InstanceSingleton) WithConfig(config config.ServerConfig) *InstanceSingleton {
	i.config = config
	return i
}

func (i *InstanceSingleton) registerMiddlewares(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost", i.config.FrontendURL, i.config.Host},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(CSPMiddleware)

	// if debug
	if config.IsDevelopment() {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.Recover())

	e.Use(middleware.Secure())

	csrfConfig := middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieDomain:   i.config.Host,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	})

	e.Use(csrfConfig)
}

func generateSecureNonce() (string, error) {
	nonce := make([]byte, config.NonceLength)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(nonce), err
}

func CSPMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		htmxNonce, _ := generateSecureNonce()
		hyperscriptNonce, _ := generateSecureNonce()
		picoCSSNonce, _ := generateSecureNonce()
		cssScopeInlineNonce, _ := generateSecureNonce()
		preloadNonce, _ := generateSecureNonce()
		umamiNonce, _ := generateSecureNonce()

		_ = "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg="

		cspHeader := fmt.Sprintf(
			"default-src 'self'; connect-src 'self' https://umami.memnix.app ; script-src 'nonce-%s' 'nonce-%s' 'nonce-%s' 'nonce-%s' 'nonce-%s'; style-src 'self' 'unsafe-inline' https://fonts.bunny.net; font-src https://fonts.bunny.net 'self'",
			htmxNonce, hyperscriptNonce, preloadNonce, cssScopeInlineNonce, umamiNonce)

		c.Response().Header().Set("Content-Security-Policy", cspHeader)

		c.Set("nonce", domain.Nonce{
			HtmxNonce:           htmxNonce,
			HyperscriptNonce:    hyperscriptNonce,
			PreloadNonce:        preloadNonce,
			UmamiNonce:          umamiNonce,
			PicoCSSNonce:        picoCSSNonce,
			CSSScopeInlineNonce: cssScopeInlineNonce,
		})

		return next(c)
	}
}

func StaticAssetsCacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
		return next(c)
	}
}

func StaticPageCacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Set Cache-Control header
		c.Response().Header().Set("Cache-Control", "private, max-age=60")

		return next(c)
	}
}
