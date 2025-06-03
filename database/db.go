package db

import (
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MustInitDB(dbUrl string, log *slog.Logger) *gorm.DB {
	gormLogger := logger.New(
		&GormLogAdapter{Logger: log},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn, // Only log slow queries and errors
			IgnoreRecordNotFoundError: true,
		},
	)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	// Configure connection pool
	dbConn, err := db.DB()
	if err != nil {
		panic("Failed to get connection pool: " + err.Error())
	}
	dbConn.SetMaxIdleConns(10)               // Number of idle connections in the pool
	dbConn.SetMaxOpenConns(100)              // Maximum number of open connections
	dbConn.SetConnMaxLifetime(1 * time.Hour) // Max time a connection can be reused

	if err := dbConn.Ping(); err != nil {
		panic("Failed to ping database: " + err.Error())
	}
	log.Info("Database connected successfully")

	return db
}

type GormLogAdapter struct {
	Logger *slog.Logger
}

func (l *GormLogAdapter) Printf(format string, args ...interface{}) {
	l.Logger.Debug(format, "args", args)
}
