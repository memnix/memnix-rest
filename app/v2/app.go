package v2

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/memnix/memnix-rest/cmd/v1/config"
)

var (
	instance *InstanceSingleton //nolint:gochecknoglobals //Singleton
	once     sync.Once          //nolint:gochecknoglobals //Singleton
)

type InstanceSingleton struct {
	echoInstance *echo.Echo
	config       ServerConfig
}

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	Port        string
	AppVersion  string
	JaegerURL   string
	Host        string
	FrontendURL string
	LogLevel    string
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

func CreateEchoInstance(config ServerConfig) *InstanceSingleton {
	return GetEchoSingleton().WithConfig(config)
}

func (i *InstanceSingleton) Start() error {
	if err := i.echoInstance.Start(":" + i.config.Port); err != nil {
		return err
	}

	return nil
}

func (i *InstanceSingleton) WithConfig(config ServerConfig) *InstanceSingleton {
	i.config = config
	return i
}

func (i *InstanceSingleton) registerMiddlewares(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost", i.config.FrontendURL, i.config.Host},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// if debug
	if config.IsDevelopment() {
		e.Use(middleware.Logger())
	}

	// e.Use(middleware.Recover())

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
