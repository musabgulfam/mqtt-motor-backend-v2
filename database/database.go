package database

import (
	"fmt"
	"log"
	"os"

	"github.com/musabgulfam/pumplink-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable that holds our database connection
// This allows other parts of our application to access the database
// We use a global variable because we only need one database connection for the entire application
var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
// This function is responsible for:
// 1. Opening a connection to the PostgreSQL database using environment variables
// 2. Setting up GORM (Go Object Relational Mapper) for database operations
// 3. Auto-migrating the database schema based on our models
// 4. Making the connection available to other parts of the application
//
// Returns:
//   - error: Any error that occurred during the connection process
func Connect() error {
	var err error

	// Build the DSN (Data Source Name) string for PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Karachi",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	) // This DSN (Data Source Name) string contains the necessary connection details
	// DSN includes:
	// - host: Database server address
	// - user: Database username
	// - password: Database user's password
	// - dbname: Name of the database to connect to
	// - port: Port number on which the database server is listening
	// DSN is a standard way to specify connection parameters for databases
	// Replace with your actual database credentials and connection details
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return err
	}
	DB = db

	// Auto-migrate the database schema based on our models
	// This ensures our database tables match our Go structs
	// GORM will create tables, add columns, or modify schema as needed
	// This is especially useful during development when models change frequently
	err = DB.AutoMigrate(
		&models.User{},
		// &models.DeviceActivationLog{},
		&models.Device{},
		&models.DeviceLog{},
		&models.DeviceSession{},
	)
	if err != nil {
		// If migration fails, return the error
		// This could happen if there are schema conflicts or database issues
		return err
	}

	// Now insert initial data (e.g., a default Motor device)
	var count int64
	DB.Model(&models.Device{}).Where("name = ?", "Motor Pump").Count(&count)
	if count == 0 {
		DB.Create(&models.Device{
			Name:  "Motor Pump",
			State: "UNKNOWN",
		})
	}

	// Return nil to indicate successful connection and migration
	return nil
}

// GetDB returns the global database connection
// This function provides a clean way for other parts of the application
// to access the database connection without directly accessing the global variable
//
// Returns:
//   - *gorm.DB: The database connection instance
func GetDB() *gorm.DB {
	return DB
}
