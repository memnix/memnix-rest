package infrastructures

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnSingleton struct {
	dbConn *gorm.DB
	config DatabaseConfig
}

var (
	dbInstance *DBConnSingleton //nolint:gochecknoglobals //Singleton
	dbOnce     sync.Once        //nolint:gochecknoglobals //Singleton
)

// DatabaseConfig holds the configuration for the database.
type DatabaseConfig struct {
	DSN             string
	SQLMaxIdleConns int
	SQLMaxOpenConns int
}

func GetDBConnInstance() *DBConnSingleton {
	dbOnce.Do(func() {
		dbInstance = &DBConnSingleton{}
	})
	return dbInstance
}

func NewDBConnInstance(config DatabaseConfig) *DBConnSingleton {
	return GetDBConnInstance().WithConfig(config)
}

func (d *DBConnSingleton) WithConfig(config DatabaseConfig) *DBConnSingleton {
	d.config = config
	return d
}

func (d *DBConnSingleton) GetDBConn() *gorm.DB {
	return d.dbConn
}

func (d *DBConnSingleton) ConnectDB() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)
	// Open connection
	conn, err := gorm.Open(postgres.Open(d.config.DSN), &gorm.Config{
		Logger:                                   newLogger, // Logger
		SkipDefaultTransaction:                   true,      // Skip default transaction
		DisableForeignKeyConstraintWhenMigrating: true,      // Disable foreign key constraint when migrating (planetscale recommends this)
	})
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}

	sqlDB, err := conn.DB() // Get sql.DB object from gorm.DB
	if err != nil {
		return errors.Wrap(err, "failed to get sql.DB object")
	}
	sqlDB.SetMaxIdleConns(d.config.SQLMaxIdleConns) // Set max idle connections
	sqlDB.SetMaxOpenConns(d.config.SQLMaxOpenConns) // Set max open connections
	sqlDB.SetConnMaxLifetime(time.Hour)             // Set connection max lifetime

	d.dbConn = conn

	return nil
}

func (d *DBConnSingleton) DisconnectDB() error {
	sqlDB, err := d.GetDBConn().DB() // Get sql.DB object from gorm.DB
	if err != nil {
		return errors.Wrap(err, "failed to get sql.DB object")
	}
	if err = sqlDB.Close(); err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}

	return nil
}

func GetDBConn() *gorm.DB {
	if GetDBConnInstance().GetDBConn() == nil {
		slog.Error("db connection is nil")
		return nil
	}
	return GetDBConnInstance().GetDBConn()
}
