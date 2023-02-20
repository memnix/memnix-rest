package infrastructures

import (
	"context"
	"github.com/rs/zerolog/log"
)

import "github.com/edgedb/edgedb-go"

var (
	edgeDBClient *edgedb.Client
)

// ConnectEdgeDB Connects to EdgeDB
func ConnectEdgeDB() error {
	edgeDBClient = CreateEdgeDBConnection()

	return nil
}

// CloseEdgeDB Closes EdgeDB connection
func CloseEdgeDB() error {
	return edgeDBClient.Close()
}

// GetEdgeDBClient Returns EdgeDB client
func GetEdgeDBClient() *edgedb.Client {
	return edgeDBClient
}

// CreateEdgeDBConnection Returns new EdgeDB client
func CreateEdgeDBConnection() *edgedb.Client {
	ctx := context.Background()
	client, err := edgedb.CreateClientDSN(ctx, "memnix_edgedb", edgedb.Options{
		TLSOptions: edgedb.TLSOptions{
			SecurityMode: edgedb.TLSModeInsecure,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create EdgeDB client")
	}

	return client

}
