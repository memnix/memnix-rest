package barista

type IUseCase interface {
	// GetName returns the name of the barista.
	GetName() string
}

type IPgRepository interface {
	// GetName returns the name of the barista.
	GetName() string
}
