package routes

import (
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/services"
)

func registerDeckRoutes() {

	deckController := services.GetServiceContainer().InjectDeckController()

	// Admin only
	routesMap["/decks"] = routeStruct{
		Method:     "GET",
		Handler:    deckController.GetAllDecks,
		Permission: models.PermAdmin,
	}

	// Get routes
	routesMap["/decks/public"] = routeStruct{
		Method:     "GET",
		Handler:    deckController.GetAllPublicDecks,
		Permission: models.PermUser,
	}

	routesMap["/decks/available"] = routeStruct{
		Method:     "GET",
		Handler:    deckController.GetAllAvailableDecks,
		Permission: models.PermUser,
	}

	// Get all editor decks
	routesMap["/decks/editor"] = routeStruct{
		Method:     "GET",
		Handler:    deckController.GetAllEditorDecks,
		Permission: models.PermUser,
	}

	// Get all decks the user is sub to
	routesMap["/decks/sub"] = routeStruct{
		Method:     "GET",
		Handler:    deckController.GetAllSubDecks,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID"] = routeStruct{
		Method:     "GET",
		Handler:    deckController.GetDeckByID,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/users"] = routeStruct{
		Method:     "GET",
		Handler:    deckController.GetAllSubUsers,
		Permission: models.PermUser,
	}

	// Post
	routesMap["/decks/new"] = routeStruct{
		Method:     "POST",
		Handler:    deckController.CreateNewDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/subscribe"] = routeStruct{
		Method:     "POST",
		Handler:    deckController.SubToDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/unsubscribe"] = routeStruct{
		Method:     "POST",
		Handler:    deckController.UnSubToDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/private/:key/:code/subscribe"] = routeStruct{
		Method:     "POST",
		Handler:    deckController.SubToPrivateDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/publish"] = routeStruct{
		Method:     "POST",
		Handler:    deckController.PublishDeckRequest,
		Permission: models.PermUser,
	}

	// Put
	routesMap["/decks/:deckID/edit"] = routeStruct{
		Method:     "PUT",
		Handler:    deckController.UpdateDeckByID,
		Permission: models.PermUser,
	}

	// Delete
	routesMap["/decks/:deckID"] = routeStruct{
		Method:     "DELETE",
		Handler:    deckController.DeleteDeckById,
		Permission: models.PermUser,
	}
}
