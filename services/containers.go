package services

import (
	"sync"

	"github.com/memnix/memnix-rest/app/v1/controllers"
	"github.com/memnix/memnix-rest/app/v2/handlers"
)

type ServiceContainer struct {
	user        controllers.UserController
	auth        controllers.AuthController
	jwt         controllers.JwtController
	oAuth       controllers.OAuthController
	deck        controllers.DeckController
	card        controllers.CardController
	mcq         controllers.McqController
	authHandler handlers.AuthController
}

var (
	container *ServiceContainer //nolint:gochecknoglobals // Singleton
	once      sync.Once         //nolint:gochecknoglobals // Singleton
)

func DefaultServiceContainer() *ServiceContainer {
	userController := InitializeUser()
	authController := InitializeAuth()
	jwtController := InitializeJWT()
	oAuthController := InitializeOAuth()
	deckController := InitializeDeck()
	cardController := InitializeCard()
	mcqController := InitializeMcq()

	authHandler := InitializeAuthHandler()

	return NewServiceContainer(userController, authController,
		jwtController, oAuthController, deckController,
		cardController, mcqController, authHandler)
}

func NewServiceContainer(user controllers.UserController, auth controllers.AuthController,
	jwt controllers.JwtController, oAuth controllers.OAuthController,
	deck controllers.DeckController, card controllers.CardController,
	mcq controllers.McqController, authHandler handlers.AuthController,
) *ServiceContainer {
	once.Do(func() {
		container = &ServiceContainer{
			user:        user,
			auth:        auth,
			jwt:         jwt,
			oAuth:       oAuth,
			deck:        deck,
			card:        card,
			mcq:         mcq,
			authHandler: authHandler,
		}
	})
	return container
}

func (sc *ServiceContainer) User() controllers.UserController {
	return sc.user
}

func (sc *ServiceContainer) Auth() controllers.AuthController {
	return sc.auth
}

func (sc *ServiceContainer) Jwt() controllers.JwtController {
	return sc.jwt
}

func (sc *ServiceContainer) OAuth() controllers.OAuthController {
	return sc.oAuth
}

func (sc *ServiceContainer) Deck() controllers.DeckController {
	return sc.deck
}

func (sc *ServiceContainer) Card() controllers.CardController {
	return sc.card
}

func (sc *ServiceContainer) Mcq() controllers.McqController {
	return sc.mcq
}

func (sc *ServiceContainer) AuthHandler() handlers.AuthController {
	return sc.authHandler
}
