package user

type IUseCase interface {
	// GetName returns the name of the user.
	GetName(id string) string
}

type IRepository interface {
	// GetName returns the name of the user.
	GetName(id uint) string
}
