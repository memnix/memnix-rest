package routes

import (
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/services"
)

func registerMcqRoutes() {
	mcqController := services.GetServiceContainer().InjectMcqController()

	routesMap["/mcqs/new"] = routeStruct{
		Method:     "POST",
		Handler:    mcqController.CreateMcq,
		Permission: models.PermUser,
	}

	routesMap["/mcqs/:deckID"] = routeStruct{
		Method:     "GET",
		Handler:    mcqController.GetMcqsByDeck,
		Permission: models.PermUser,
	}

	routesMap["/mcqs/:id/edit"] = routeStruct{
		Method:     "PUT",
		Handler:    mcqController.UpdateMcqByID,
		Permission: models.PermUser,
	}

	routesMap["/mcqs/:id"] = routeStruct{
		Method:     "DELETE",
		Handler:    mcqController.DeleteMcqByID,
		Permission: models.PermUser,
	}

}
