package services

import (
	"sync"

	"github.com/memnix/memnix-rest/app/v2/handlers"
)

type ServiceContainer struct {
	authHandler handlers.AuthController
}

var (
	container *ServiceContainer //nolint:gochecknoglobals // Singleton
	once      sync.Once         //nolint:gochecknoglobals // Singleton
)

func DefaultServiceContainer() *ServiceContainer {
	authHandler := InitializeAuthHandler()

	return NewServiceContainer(authHandler)
}

func NewServiceContainer(authHandler handlers.AuthController,
) *ServiceContainer {
	once.Do(func() {
		container = &ServiceContainer{
			authHandler: authHandler,
		}
	})
	return container
}

func (sc *ServiceContainer) AuthHandler() handlers.AuthController {
	return sc.authHandler
}
