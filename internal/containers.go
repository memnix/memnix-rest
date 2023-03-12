package internal

import (
	"sync"

	"github.com/memnix/memnix-rest/app/http/controllers"
)

type ServiceContainer interface {
	GetUser() controllers.UserController
	GetAuth() controllers.AuthController
	GetJwt() controllers.JwtController
	GetOAuth() controllers.OAuthController
	GetDeck() controllers.DeckController
}

type kernel struct{}

func (kernel) GetUser() controllers.UserController {
	return InitializeUser()
}

func (kernel) GetAuth() controllers.AuthController {
	return InitializeAuth()
}

func (kernel) GetJwt() controllers.JwtController {
	return InitializeJWT()
}

func (kernel) GetOAuth() controllers.OAuthController {
	return InitializeOAuth()
}

func (kernel) GetDeck() controllers.DeckController {
	return InitializeDeck()
}

var (
	k             *kernel   // kernel is the service container
	containerOnce sync.Once // containerOnce is used to ensure that the service container is only initialized once
)

// GetServiceContainer returns the service container
func GetServiceContainer() ServiceContainer {
	containerOnce.Do(func() {
		k = &kernel{}
	})
	return k
}
