package services

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"github.com/memnix/memnixrest/interfaces"
)

type DeckService struct {
	interfaces.IDeckService
}

func (k *kernel) InjectDeckController() controllers.DeckController {
	DBConn := infrastructures.GetDBConn()

	deckRepository := &repositories.DeckRepository{DBConn: DBConn}
	deckService := &DeckService{IDeckService: deckRepository}
	deckController := controllers.DeckController{IDeckService: deckService}

	return deckController
}
