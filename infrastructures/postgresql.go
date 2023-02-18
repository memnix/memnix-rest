package infrastructures

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

// DBParams is the database configuration parameters
type DBParams struct {
	Host     string // Database host
	Port     string // Database port
	User     string // Database user
	Password string // Database password
	DBName   string // Database name
}

var (
	// DBConn is the database connection
	dbConn *gorm.DB
)

// GetDBConn returns the database connection object
func GetDBConn() *gorm.DB {
	return dbConn
}

// DBConfig returns the database configuration parameters
// from environment variables
func DBConfig() DBParams {
	dbConfig := new(DBParams) // Create a new DBParams struct

	// Get database configuration from environment variables
	dbConfig.User = os.Getenv("DB_USER")         // Get DB_USER from env
	dbConfig.Password = os.Getenv("DB_PASSWORD") // Get DB_PASSWORD from env
	dbConfig.Host = os.Getenv("DB_HOST")         // Get DB_HOST from env
	dbConfig.DBName = os.Getenv("DB_DB")         // Get DB_DB (db name) from env
	dbConfig.Port = os.Getenv("DB_PORT")         // Get DB_PORT from env

	return *dbConfig
}

func (db *DBParams) GetConnectionString() string {
	dbPort, _ := strconv.Atoi(db.Port)

	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable",
		db.User, db.Password, db.Host, db.DBName, dbPort)
}

// ConnectDB connects to the database
func ConnectDB() error {
	dbConfig := DBConfig() // Get database configuration parameters

	// Create a new database connection
	db, err := gorm.Open(postgres.Open(dbConfig.GetConnectionString()), &gorm.Config{})
	if err != nil {
		return err
	}

	// Set the database connection
	dbConn = db

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
