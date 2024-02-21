package user

import (
	"context"

	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/utils"
)

// UseCase is the user use case.
type UseCase struct {
	IRepository
	IRedisRepository
	IRistrettoRepository
}

// GetName returns the name of the user with the given id.
func (u UseCase) GetName(ctx context.Context, id string) string {
	uintID, _ := utils.ConvertStrToUInt(id)

	return u.IRepository.GetName(ctx, uintID)
}

// GetByID returns the user with the given id.
func (u UseCase) GetByID(ctx context.Context, id uint) (domain.User, error) {
	_, span := infrastructures.GetTracerInstance().Tracer().Start(ctx, "GetUserByID")
	defer span.End()

	var userObject domain.User

	if risrettoHit, err := u.IRistrettoRepository.Get(ctx, id); err == nil {
		return risrettoHit, nil
	}

	userObject, err := u.IRepository.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	_ = u.IRistrettoRepository.Set(ctx, userObject)

	return userObject, nil
}

// NewUseCase returns a new user use case.
func NewUseCase(repo IRepository, redis IRedisRepository, ristretto IRistrettoRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis, IRistrettoRepository: ristretto}
}
