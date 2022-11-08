package infrastructures

import (
	"fmt"
	"github.com/memnix/memnixrest/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

type dbParams struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var (
	// DBConn is the infrastructures connection
	dbConn *gorm.DB
)

func GetDBConn() *gorm.DB {
	return dbConn
}

// Connect creates a connection to infrastructures
func Connect() error {
	dbConfig := loadDBVar() // Load var from .env file

	// Convert port
	dbPort, err := strconv.Atoi(dbConfig.Port)
	if err != nil {
		return err
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
		},
	)

	// Create postgres connection string
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.DBName, dbPort)
	// Open connection
	dbConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(utils.SQLMaxIdleConns)
	sqlDB.SetMaxOpenConns(utils.SQLMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
