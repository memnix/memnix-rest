package infrastructures

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/memnix/memnix-rest/config"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnSingleton struct {
	dbConn *gorm.DB
}

var (
	dbInstance *DBConnSingleton //nolint:gochecknoglobals //Singleton
	dbOnce     sync.Once        //nolint:gochecknoglobals //Singleton
)

func GetDBConnInstance() *DBConnSingleton {
	dbOnce.Do(func() {
		dbInstance = &DBConnSingleton{}
	})
	return dbInstance
}

func (d *DBConnSingleton) GetDBConn() *gorm.DB {
	return d.dbConn
}

func (d *DBConnSingleton) ConnectDB(dsn string) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)
	// Open connection
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
	sqlDB.SetMaxIdleConns(config.SQLMaxIdleConns) // Set max idle connections
	sqlDB.SetMaxOpenConns(config.SQLMaxOpenConns) // Set max open connections
	sqlDB.SetConnMaxLifetime(time.Second)         // Set max connection lifetime

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
	return GetDBConnInstance().GetDBConn()
}
