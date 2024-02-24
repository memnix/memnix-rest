package v2

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/memnix/memnix-rest/app/v2/handlers"
	"github.com/memnix/memnix-rest/pkg/i18n"
	"github.com/memnix/memnix-rest/services"
)

func (i *InstanceSingleton) registerStaticRoutes(e *echo.Echo) {
	e.Static("/assets/", "assets/static")
	e.Static("/img", "assets/img")
}

func (i *InstanceSingleton) registerRoutes(e *echo.Echo) {
	serviceContainer := services.DefaultServiceContainer()
	authController := serviceContainer.AuthHandler()
	pageController := handlers.NewPageController()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("lang", "en")
			c.Set("localizer", i18n.DefaultLocalizer())

			return next(c)
		}
	})

	e.GET("/", pageController.GetIndex)
	e.GET("/login", pageController.GetLogin)
	e.GET("/register", pageController.GetRegister)
	e.POST("/register", authController.PostRegister)
	e.POST("/logout", authController.PostLogout)
	e.POST("/login", authController.PostLogin)
	e.POST("/clicked", pageController.PostClicked)

	g := e.Group(":lang")

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Debug("🌐 Language middleware")
			lang := c.Param("lang")
			if lang != "en" && lang != "fr" {
				lang = "en"
			}

			c.Set("lang", lang)
			accept := c.Request().Header.Get("Accept-Language")

			localizer := i18n.Localizer(lang, accept)
			if localizer == nil {
				localizer = i18n.DefaultLocalizer()
			}

			c.Set("localizer", localizer)

			return next(c)
		}
	})

	g.Add("GET", "/", pageController.GetIndex)
}
