package user

import (
	"context"
	"github.com/edgedb/edgedb-go"
	"github.com/rs/zerolog/log"
)

type EdgeRepository struct {
	EdgedbConn *edgedb.Client
}

func NewEdgeRepository(edgedbConn *edgedb.Client) IEdgeRepository {
	return &EdgeRepository{EdgedbConn: edgedbConn}
}

// GetName returns the name of the user.
func (r *EdgeRepository) GetName(id edgedb.UUID) string {
	query := `SELECT User { username } FILTER .id = <uuid>$id`
	var user struct {
		Username string `edgedb:"username"`
	}

	err := r.EdgedbConn.QuerySingle(
		context.Background(), query, &user, map[string]interface{}{"id": id})
	if err != nil {
		log.Debug().Err(err).Msg("Error getting user")
	}

	return user.Username
}
