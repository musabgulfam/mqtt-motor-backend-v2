package database

import (
	"mqtt-motor-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a global variable that holds our database connection
// This allows other parts of our application to access the database
// We use a global variable because we only need one database connection for the entire application
var DB *gorm.DB

// Connect establishes a connection to the SQLite database
// This function is responsible for:
// 1. Opening a connection to the SQLite database file
// 2. Setting up GORM (Go Object Relational Mapper) for database operations
// 3. Creating the database file if it doesn't exist
// 4. Auto-migrating the database schema based on our models
// 5. Making the connection available to other parts of the application
//
// Parameters:
//   - dbPath: The path to the SQLite database file (e.g., "data.db")
//
// Returns:
//   - error: Any error that occurred during the connection process
func Connect(dbPath string) error {
	var err error

	// Open a connection to the SQLite database using GORM
	// sqlite.Open() creates a new SQLite driver instance
	// &gorm.Config{} provides default GORM configuration
	// If the database file doesn't exist, SQLite will create it automatically
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// If connection fails, return the error
		// Common causes: insufficient permissions, disk full, invalid path
		return err
	}

	// Auto-migrate the database schema based on our models
	// This ensures our database tables match our Go structs
	// GORM will create tables, add columns, or modify schema as needed
	// This is especially useful during development when models change frequently
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		// If migration fails, return the error
		// This could happen if there are schema conflicts or database issues
		return err
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
