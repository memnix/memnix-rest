package routes

import (
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/services"
)

func registerUserRoutes() {
	u := services.GetServiceContainer().InjectUserController()

	// Get
	routesMap["/users"] = routeStruct{
		Method:     "GET",
		Handler:    u.GetAllUsers,
		Permission: models.PermAdmin,
	}

	routesMap["/users/id/:id"] = routeStruct{
		Method:     "GET",
		Handler:    u.GetUserByID,
		Permission: models.PermAdmin,
	}

	// Post
	routesMap["/users/settings/:deckID/today"] = routeStruct{
		Method:     "POST",
		Handler:    u.SetTodayConfig,
		Permission: models.PermUser,
	}

	routesMap["/users/settings/resetpassword"] = routeStruct{
		Method:     "POST",
		Handler:    u.ResetPassword,
		Permission: models.PermUser,
	}

	routesMap["/users/settings/confirmpassword"] = routeStruct{
		Method:     "POST",
		Handler:    u.ResetPasswordConfirm,
		Permission: models.PermUser,
	}

	// Put
	routesMap["/users/id/:id"] = routeStruct{
		Method:     "PUT",
		Handler:    u.UpdateUserByID,
		Permission: models.PermAdmin,
	}

}
