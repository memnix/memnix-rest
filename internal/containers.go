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
}

type kernel struct{}

func (k kernel) GetUser() controllers.UserController {
	return InitializeUser()
}

func (k kernel) GetAuth() controllers.AuthController {
	return InitializeAuth()
}

func (k kernel) GetJwt() controllers.JwtController {
	return InitializeJWT()
}

func (k kernel) GetOAuth() controllers.OAuthController {
	return InitializeOAuth()
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
