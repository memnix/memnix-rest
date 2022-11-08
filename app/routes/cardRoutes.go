package routes

import (
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/services"
)

func registerCardRoutes() {

	c := services.GetServiceContainer().InjectCardController()

	// Get routes
	routesMap["/cards/today"] = routeStruct{
		Method:     "GET",
		Handler:    c.GetAllTodayCard,
		Permission: models.PermUser,
	}

	routesMap["/cards/:deckID/training"] = routeStruct{
		Method:     "GET",
		Handler:    c.GetTrainingCardsByDeck,
		Permission: models.PermUser,
	}

	routesMap["/cards/deck/:deckID"] = routeStruct{
		Method:     "GET",
		Handler:    c.GetCardsFromDeck,
		Permission: models.PermUser,
	}

	routesMap["/cards/new"] = routeStruct{
		Method:     "POST",
		Handler:    c.CreateNewCard,
		Permission: models.PermUser,
	}

	// Post routes
	routesMap["/cards/response"] = routeStruct{
		Method:     "POST",
		Handler:    c.PostResponse,
		Permission: models.PermUser,
	}

	routesMap["/cards/selfresponse"] = routeStruct{
		Method:     "POST",
		Handler:    c.PostSelfEvaluateResponse,
		Permission: models.PermUser,
	}

	// Put routes

	routesMap["/cards/:id/edit"] = routeStruct{
		Method:     "PUT",
		Handler:    c.UpdateCardByID,
		Permission: models.PermUser,
	}

	// Delete routes

	routesMap["/cards/:id"] = routeStruct{
		Method:     "DELETE",
		Handler:    c.DeleteCardByID,
		Permission: models.PermUser,
	}

	// ADMIN ONLY

	routesMap["/cards"] = routeStruct{
		Method:     "GET",
		Handler:    c.GetAllCards,
		Permission: models.PermAdmin,
	}

	routesMap["/cards/id/:id"] = routeStruct{
		Method:     "GET",
		Handler:    c.GetCardByID,
		Permission: models.PermAdmin,
	}

}
