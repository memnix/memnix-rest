package services

import (
	"github.com/memnix/memnixrest/app/controllers"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/data/repositories"
	"github.com/memnix/memnixrest/interfaces"
)

type CardService struct {
	interfaces.ICardRepository
}

func (k *kernel) InjectCardController() controllers.CardController {
	DBConn := infrastructures.GetDBConn()

	cardRepository := &repositories.CardRepository{DBConn: DBConn}
	cardService := &CardService{ICardRepository: cardRepository}
	cardController := controllers.CardController{ICardService: cardService}

	return cardController
}
