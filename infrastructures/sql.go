package infrastructures

import (
	"log"
	"os"
	"time"

	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/env"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
)

// DBConn is the database connection object
var DBConn *gorm.DB

// GetDBConn returns the database connection object
func GetDBConn() *gorm.DB {
	return DBConn
}

// getDSN returns the database connection string
// see: utils/env.go
func getDSN(env *env.Env) string {
	var dsn string

	// Get database configuration from environment variables
	if config.IsDevelopment() {
		dsn = env.GetEnv("DEBUG_DB_DSN")
	} else {
		dsn = env.GetEnv("DB_DSN")
	}
	return dsn
}

// ConnectDB creates a connection to database
//
// see: utils/config.go and utils/env.go for more details
func ConnectDB() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	dsn := getDSN(config.EnvHelper)

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
	sqlDB.SetConnMaxLifetime(time.Hour)           // Set max connection lifetime

	if err = conn.Use(prometheus.New(prometheus.Config{
		DBName:          "db1",                                // `DBName` as metrics label
		RefreshInterval: config.GormPrometheusRefreshInterval, // refresh metrics interval (default 15 seconds)
		StartServer:     false,                                // start http server to expose metrics
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{VariableNames: []string{"Threads_running"}},
		},
	})); err != nil {
		return err
	}

	if err = conn.Use(tracing.NewPlugin()); err != nil {
		return err
	}

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
