package services

import (
	"github.com/memnix/memnixrest/app/controllers"
	"sync"
)

type IServiceContainer interface {
	InjectUserController() controllers.UserController
	InjectAuthController() controllers.AuthController
	InjectCardController() controllers.CardController
	InjectDeckController() controllers.DeckController
	InjectMcqController() controllers.McqController
}

type kernel struct{}

var (
	k             *kernel
	containerOnce sync.Once
)

func GetServiceContainer() IServiceContainer {
	containerOnce.Do(func() {
		k = &kernel{}
	})
	return k
}
