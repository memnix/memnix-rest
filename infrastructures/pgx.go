package infrastructures

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConnSingleton struct {
	conn   *pgxpool.Pool
	config PgxConfig
}

var (
	pgxInstance *PgxConnSingleton //nolint:gochecknoglobals //Singleton
	pgxOnce     sync.Once         //nolint:gochecknoglobals //Singleton
)

type PgxConfig struct {
	DSN string
}

func GetPgxConnInstance() *PgxConnSingleton {
	pgxOnce.Do(func() {
		pgxInstance = &PgxConnSingleton{}
	})
	return pgxInstance
}

func NewPgxConnInstance(config PgxConfig) *PgxConnSingleton {
	return GetPgxConnInstance().WithConfig(config)
}

func (p *PgxConnSingleton) WithConfig(config PgxConfig) *PgxConnSingleton {
	p.config = config
	return p
}

func (p *PgxConnSingleton) GetPgxConn() *pgxpool.Pool {
	return p.conn
}

func (p *PgxConnSingleton) ConnectPgx() error {
	conn, err := pgxpool.New(context.Background(), p.config.DSN)
	if err != nil {
		return err
	}

	p.conn = conn
	return nil
}

func (p *PgxConnSingleton) ClosePgx() {
	p.conn.Close()
}

func GetPgxConn() *pgxpool.Pool {
	return pgxInstance.GetPgxConn()
}
