package routes

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/pkg/models"
)

func registerUserRoutes() {
	// Get
	routesMap["/users"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllUsers,
		Permission: models.PermAdmin,
	}

	routesMap["/users/id/:id"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetUserByID,
		Permission: models.PermAdmin,
	}

	// Post
	routesMap["/users/settings/:deckID/today"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.SetTodayConfig,
		Permission: models.PermUser,
	}

	routesMap["/users/settings/resetpassword"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.ResetPassword,
		Permission: models.PermUser,
	}

	routesMap["/users/settings/confirmpassword"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.ResetPasswordConfirm,
		Permission: models.PermUser,
	}

	// Put
	routesMap["/users/id/:id"] = routeStruct{
		Method:     "PUT",
		Handler:    controllers.UpdateUserByID,
		Permission: models.PermAdmin,
	}

}
