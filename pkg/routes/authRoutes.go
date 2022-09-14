package routes

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/app/models"
)

func registerAuthRoutes() {
	// Get
	routesMap["/user"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.User,
		Permission: models.PermUser,
	}

	// Post
	routesMap["/login"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.Login,
		Permission: models.PermNone,
	}

	routesMap["/register"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.Register,
		Permission: models.PermNone,
	}

	routesMap["/logout"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.Logout,
		Permission: models.PermNone,
	}
}
