package user

import (
	"context"

	"github.com/dgraph-io/ristretto"
	db "github.com/memnix/memnix-rest/db/sqlc"
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

func (r *RistrettoRepository) Get(ctx context.Context, id int32) (db.User, error) {
	_, span := infrastructures.GetTracerInstance().Tracer().Start(ctx, "GetByIDRistretto")
	defer span.End()

	ristrettoHit, ok := r.RistrettoCache.Get(keyPrefix + utils.ConvertInt32ToStr(id))
	if !ok {
		return db.User{}, errors.New("user not found")
	}

	switch ristrettoHit := ristrettoHit.(type) {
	case db.User:
		return ristrettoHit, nil
	default:
		return db.User{}, errors.New("user not found")
	}
}

func (r *RistrettoRepository) Set(ctx context.Context, user db.User) error {
	_, span := infrastructures.GetTracerInstance().Tracer().Start(ctx, "SetByIDRistretto")
	defer span.End()

	r.RistrettoCache.Set(keyPrefix+utils.ConvertInt32ToStr(user.ID), user, 0)

	r.RistrettoCache.Wait()

	return nil
}

func (r *RistrettoRepository) Delete(ctx context.Context, id int32) error {
	_, span := infrastructures.GetTracerInstance().Tracer().Start(ctx, "DeleteByIDRistretto")
	defer span.End()

	r.RistrettoCache.Del(keyPrefix + utils.ConvertInt32ToStr(id))

	return nil
}
