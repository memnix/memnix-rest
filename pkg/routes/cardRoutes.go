package routes

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/app/models"
)

func registerCardRoutes() {

	// Get routes
	routesMap["/cards/today"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllTodayCard,
		Permission: models.PermUser,
	}

	routesMap["/cards/:deckID/training"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetTrainingCardsByDeck,
		Permission: models.PermUser,
	}

	routesMap["/mcqs/:deckID"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetMcqsByDeck,
		Permission: models.PermUser,
	}

	routesMap["/cards/deck/:deckID"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetCardsFromDeck,
		Permission: models.PermUser,
	}

	routesMap["/cards/new"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.CreateNewCard,
		Permission: models.PermUser,
	}

	routesMap["/mcqs/new"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.CreateMcq,
		Permission: models.PermUser,
	}

	// Post routes
	routesMap["/cards/response"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.PostResponse,
		Permission: models.PermUser,
	}

	routesMap["/cards/selfresponse"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.PostSelfEvaluateResponse,
		Permission: models.PermUser,
	}

	// Put routes

	routesMap["/cards/:id/edit"] = routeStruct{
		Method:     "PUT",
		Handler:    controllers.UpdateCardByID,
		Permission: models.PermUser,
	}

	routesMap["/mcqs/:id/edit"] = routeStruct{
		Method:     "PUT",
		Handler:    controllers.UpdateMcqByID,
		Permission: models.PermUser,
	}

	// Delete routes

	routesMap["/cards/:id"] = routeStruct{
		Method:     "DELETE",
		Handler:    controllers.DeleteCardByID,
		Permission: models.PermUser,
	}

	routesMap["/mcqs/:id"] = routeStruct{
		Method:     "DELETE",
		Handler:    controllers.DeleteMcqByID,
		Permission: models.PermUser,
	}

	// ADMIN ONLY

	routesMap["/cards"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllCards,
		Permission: models.PermAdmin,
	}

	routesMap["/cards/id/:id"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetCardByID,
		Permission: models.PermAdmin,
	}

}
