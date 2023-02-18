package barista

type UseCase struct {
	IPgRepository
}

func NewUseCase(pgRepository IPgRepository) IUseCase {
	return &UseCase{IPgRepository: pgRepository}
}

// GetName returns the name of the barista.
func (u *UseCase) GetName() string {
	return u.IPgRepository.GetName()
}
