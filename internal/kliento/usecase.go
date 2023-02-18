package kliento

type UseCase struct {
	IRedisRepository
}

func (u UseCase) GetName() string {
	return u.IRedisRepository.GetName()
}

func (u UseCase) SetName(name string) error {
	return u.IRedisRepository.SetName(name)
}

func NewUseCase(redisRepo IRedisRepository) IUseCase {
	return &UseCase{IRedisRepository: redisRepo}
}
