package user

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/utils"
	"github.com/pkg/errors"
)

type RistrettoRepository struct {
	RistrettoCache *ristretto.Cache // RistrettoCache is the cache
}

func NewRistrettoCache(ristrettoCache *ristretto.Cache) IRistrettoRepository {
	return &RistrettoRepository{
		RistrettoCache: ristrettoCache,
	}
}

func (r *RistrettoRepository) Get(ctx context.Context, id uint) (domain.User, error) {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "GetByIDRistretto")
	defer span.End()

	ristrettoHit, ok := r.RistrettoCache.Get("user:" + utils.ConvertUIntToStr(id))
	if !ok {
		return domain.User{}, errors.New("user not found")
	}

	switch ristrettoHit.(type) {
	case domain.User:
		span.AddEvent("got user" + utils.ConvertUIntToStr(ristrettoHit.(domain.User).ID) + " " + ristrettoHit.(domain.User).Email)
		return ristrettoHit.(domain.User), nil
	default:
		return domain.User{}, errors.New("user not found")
	}
}

func (r *RistrettoRepository) Set(ctx context.Context, user domain.User) error {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "SetByIDRistretto")
	defer span.End()

	r.RistrettoCache.Set("user:"+utils.ConvertUIntToStr(user.ID), user, 0)

	r.RistrettoCache.Wait()

	return nil
}

func (r *RistrettoRepository) Delete(ctx context.Context, id uint) error {
	_, span := infrastructures.GetFiberTracer().Start(ctx, "DeleteByIDRistretto")
	defer span.End()

	r.RistrettoCache.Del("user:" + utils.ConvertUIntToStr(id))

	return nil
}
