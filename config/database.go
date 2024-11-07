package config

import (
	"fmt"
	"log"

	"github.com/wanliqun/go-wallet-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Auto migrating table models
var allModels = []interface{}{
	&models.User{},
	&models.Vault{},
	&models.Transaction{},
}

type DatabaseConfig struct {
	Host     string `default:"127.0.0.1"`
	Port     string `default:"5432"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	Database string `default:"wallet_db"`
	SSLMode  string `default:"disable"`
}

// MustOpenOrCreate creates an instance of store or panics on any error.
func (config *DatabaseConfig) MustOpenOrCreate() *gorm.DB {
	// Create the database if absent and return if it was newly created
	newCreated := config.mustCreateDatabaseIfAbsent()

	// Connect to the specified database
	db := config.mustConnect(config.Database)

	// Auto-migrate tables if the database was newly created
	if newCreated {
		config.autoMigrateTables(db)
	}

	log.Println("PostgreSQL database initialized")
	return db
}

// mustConnect creates a new database connection or panics on error.
func (config *DatabaseConfig) mustConnect(database string) *gorm.DB {
	dsn := config.buildDSN(database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database (%s): %v", database, err)
	}
	return db
}

// buildDSN constructs the PostgreSQL DSN string.
func (config *DatabaseConfig) buildDSN(database string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, database, config.Port, config.SSLMode,
	)
}

// mustCreateDatabaseIfAbsent checks if the database exists, and creates it if it doesn’t.
func (config *DatabaseConfig) mustCreateDatabaseIfAbsent() bool {
	// Connect to the PostgreSQL server without specifying a specific database
	db := config.mustConnect("postgres")
	defer config.closeDBConnection(db)

	// Check if the database exists
	var exists bool
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", config.Database).Scan(&exists).Error
	if err != nil {
		log.Fatalf("failed to check database existence: %v", err)
	}

	// If the database exists, return false
	if exists {
		return false
	}

	// Create the database if it does not exist
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Database)).Error; err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	log.Println("Database created for the first time")
	return true
}

// autoMigrateTables performs auto migration for the specified models.
func (config *DatabaseConfig) autoMigrateTables(db *gorm.DB) {
	if err := db.AutoMigrate(allModels...); err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}
}

// closeDBConnection safely closes the database connection.
func (config *DatabaseConfig) closeDBConnection(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err == nil {
		sqlDb.Close()
	}
}
