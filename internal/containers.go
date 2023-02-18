package internal

import (
	"github.com/memnix/memnix-rest/app/http/controllers"
	"sync"
)

type ServiceContainer interface {
	GetBarista() controllers.BaristaController
	GetKliento() controllers.KlientoController
}

type kernel struct{}

func (k kernel) GetBarista() controllers.BaristaController {
	return InitializeBarista()
}

func (k kernel) GetKliento() controllers.KlientoController {
	return InitializeKliento()
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
