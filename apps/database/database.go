package database

import (
	"chatapp-api/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a global variable for database connection
var DB *gorm.DB

// ConnectDatabase to open the connection to the PostgreSQL database
func ConnectDatabase(config *config.Config) *gorm.DB{
	var err error

	// Make DSN (Data Source Name) for PostgreSQL connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.SSLMode,
	)

	// Decide log level based on environment
	var logLevel logger.LogLevel
	if config.App.Env == "development" {
		logLevel = logger.Info // Show all query on development
	} else {
		logLevel = logger.Error // Only show error on production
	}

	// Open the connection to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")
	return DB
}