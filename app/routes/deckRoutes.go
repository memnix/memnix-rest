package routes

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/models"
)

func registerDeckRoutes() {

	// Admin only
	routesMap["/decks"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllDecks,
		Permission: models.PermAdmin,
	}

	// Get routes
	routesMap["/decks/public"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllPublicDecks,
		Permission: models.PermUser,
	}

	routesMap["/decks/available"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllAvailableDecks,
		Permission: models.PermUser,
	}

	// Get all editor decks
	routesMap["/decks/editor"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllEditorDecks,
		Permission: models.PermUser,
	}

	// Get all decks the user is sub to
	routesMap["/decks/sub"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllSubDecks,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetDeckByID,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/users"] = routeStruct{
		Method:     "GET",
		Handler:    controllers.GetAllSubUsers,
		Permission: models.PermUser,
	}

	// Post
	routesMap["/decks/new"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.CreateNewDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/subscribe"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.SubToDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/unsubscribe"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.UnSubToDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/private/:key/:code/subscribe"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.SubToPrivateDeck,
		Permission: models.PermUser,
	}

	routesMap["/decks/:deckID/publish"] = routeStruct{
		Method:     "POST",
		Handler:    controllers.PublishDeckRequest,
		Permission: models.PermUser,
	}

	// Put
	routesMap["/decks/:deckID/edit"] = routeStruct{
		Method:     "PUT",
		Handler:    controllers.UpdateDeckByID,
		Permission: models.PermUser,
	}

	// Delete
	routesMap["/decks/:deckID"] = routeStruct{
		Method:     "DELETE",
		Handler:    controllers.DeleteDeckById,
		Permission: models.PermUser,
	}
}
