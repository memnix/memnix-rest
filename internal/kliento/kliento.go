package kliento

type IUseCase interface {
	// GetName returns the name of the kliento.
	GetName() string
	// SetName sets the name of the kliento.
	SetName(name string) error
}

type IRedisRepository interface {
	// GetName returns the name of the kliento.
	GetName() string
	// SetName sets the name of the kliento.
	SetName(name string) error
}
