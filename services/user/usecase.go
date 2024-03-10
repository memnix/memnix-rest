package user

import (
	"context"

	db "github.com/memnix/memnix-rest/db/sqlc"
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
	idInt32, err := utils.ConvertStrToInt32(id)
	if err != nil {
		return ""
	}
	return u.IRepository.GetName(ctx, idInt32)
}

// GetByID returns the user with the given id.
func (u UseCase) GetByID(ctx context.Context, id int32) (db.User, error) {
	_, span := infrastructures.GetTracerInstance().Tracer().Start(ctx, "GetUserByID")
	defer span.End()

	if risrettoHit, err := u.IRistrettoRepository.Get(ctx, id); err == nil {
		return risrettoHit, nil
	}

	userObject, err := u.IRepository.GetByID(ctx, id)
	if err != nil {
		return db.User{}, err
	}

	_ = u.IRistrettoRepository.Set(ctx, userObject)

	return userObject, nil
}

// NewUseCase returns a new user use case.
func NewUseCase(repo IRepository, redis IRedisRepository, ristretto IRistrettoRepository) IUseCase {
	return &UseCase{IRepository: repo, IRedisRepository: redis, IRistrettoRepository: ristretto}
}
