package routes

import (
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/services"
)

func registerAuthRoutes() {
	a := services.GetServiceContainer().InjectAuthController()

	// Get
	routesMap["/user"] = routeStruct{
		Method:     "GET",
		Handler:    a.User,
		Permission: models.PermUser,
	}

	// Post
	routesMap["/login"] = routeStruct{
		Method:     "POST",
		Handler:    a.Login,
		Permission: models.PermNone,
	}

	routesMap["/register"] = routeStruct{
		Method:     "POST",
		Handler:    a.Register,
		Permission: models.PermNone,
	}

	routesMap["/logout"] = routeStruct{
		Method:     "POST",
		Handler:    a.Logout,
		Permission: models.PermNone,
	}
}
