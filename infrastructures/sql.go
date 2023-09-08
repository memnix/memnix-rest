package infrastructures

import (
	"log"
	"os"
	"time"

	"github.com/memnix/memnix-rest/config"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConn is the database connection object
var DBConn *gorm.DB

// GetDBConn returns the database connection object
func GetDBConn() *gorm.DB {
	return DBConn
}

// ConnectDB creates a connection to database
//
// see: utils/config.go and utils/env.go for more details
func ConnectDB(dsn string) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)
	// Open connection
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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

	DBConn = conn

	return nil
}

// DisconnectDB closes the database connection
func DisconnectDB() error {
	sqlDB, err := GetDBConn().DB() // Get sql.DB object from gorm.DB
	if err != nil {
		return errors.Wrap(err, "failed to get sql.DB object")
	}
	if err = sqlDB.Close(); err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}

	return nil
}
