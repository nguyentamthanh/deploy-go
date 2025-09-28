package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB global database connection
var DB *gorm.DB

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetDatabaseConfig returns database configuration from environment variables
func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "deploy_go"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ConnectDatabase establishes connection to PostgreSQL database using GORM
func ConnectDatabase() error {
	config := GetDatabaseConfig()

	// Construct DSN (Data Source Name) for GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	// Open connection with GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	DB = db

	// Auto migrate database tables
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	fmt.Println("✅ Database connected successfully with GORM")
	return nil
}

// CloseDatabase closes database connection
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// autoMigrate automatically creates and migrates database tables using GORM
func autoMigrate() error {
	// Auto migrate User and Post models
	if err := DB.AutoMigrate(&User{}, &Post{}); err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	fmt.Println("✅ Database tables auto-migrated successfully")
	return nil
}
