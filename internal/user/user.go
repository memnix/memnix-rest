package user

import "github.com/edgedb/edgedb-go"

type IUseCase interface {
	// GetName returns the name of the user.
	GetName(id edgedb.UUID) string
}

type IEdgeRepository interface {
	// GetName returns the name of the user.
	GetName(id edgedb.UUID) string
}
