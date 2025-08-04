package database

import (
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
// 4. Making the connection available to other parts of the application
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

	// Note: Auto-migration will be added when we create models in Phase 2
	// Auto-migration automatically creates database tables based on our Go structs
	// This ensures our database schema matches our application models
	// Example: err = DB.AutoMigrate(&models.User{}, &models.DeviceActivation{})

	// Return nil to indicate successful connection
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
