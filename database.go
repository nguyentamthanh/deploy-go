package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

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
	// Check if DATABASE_URL is provided (Render's default)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		fmt.Printf("🔗 Using DATABASE_URL for connection\n")
		return parseDatabaseURL(databaseURL)
	}

	// Fallback to individual environment variables
	fmt.Printf("🔗 Using individual environment variables for connection\n")
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

// parseDatabaseURL parses DATABASE_URL environment variable
func parseDatabaseURL(databaseURL string) *DatabaseConfig {
	u, err := url.Parse(databaseURL)
	if err != nil {
		fmt.Printf("⚠️ Failed to parse DATABASE_URL: %v, falling back to individual env vars\n", err)
		return &DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "deploy_go"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		}
	}

	// Extract database name from path
	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		dbName = "deploy_go"
	}

	// Determine SSL mode based on scheme
	sslMode := "require"
	if u.Scheme == "postgres" {
		sslMode = "require"
	} else if u.Scheme == "postgresql" {
		sslMode = "disable"
	}

	// Extract password from URL
	password, _ := u.User.Password()

	// Set default port if not provided
	port := u.Port()
	if port == "" {
		port = "5432"
	}

	return &DatabaseConfig{
		Host:     u.Hostname(),
		Port:     port,
		User:     u.User.Username(),
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}

// ConnectDatabase establishes connection to PostgreSQL database using GORM
func ConnectDatabase() error {
	// Try DATABASE_URL first (Render's preferred method)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		fmt.Printf("🔗 Using DATABASE_URL for connection\n")
		return connectWithDatabaseURL(databaseURL)
	}

	// Fallback to individual environment variables
	fmt.Printf("🔗 Using individual environment variables for connection\n")
	config := GetDatabaseConfig()

	// Construct DSN (Data Source Name) for GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	// Log connection details (without password for security)
	fmt.Printf("🔗 Connecting to database: host=%s user=%s dbname=%s port=%s sslmode=%s\n",
		config.Host, config.User, config.DBName, config.Port, config.SSLMode)

	// Open connection with GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Printf("❌ Database connection failed: %v\n", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return setupDatabase(db)
}

// connectWithDatabaseURL connects using DATABASE_URL
func connectWithDatabaseURL(databaseURL string) error {
	// Open connection with GORM using DATABASE_URL directly
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Printf("❌ Database connection with DATABASE_URL failed: %v\n", err)
		return fmt.Errorf("failed to connect to database with DATABASE_URL: %w", err)
	}

	return setupDatabase(db)
}

// setupDatabase configures the database connection
func setupDatabase(db *gorm.DB) error {
	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	DB = db

	// Test connection with retry mechanism
	if err := testConnection(db); err != nil {
		return fmt.Errorf("failed to test database connection: %w", err)
	}

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

// testConnection tests the database connection with retry mechanism
func testConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Retry connection up to 5 times with 2 second intervals
	for i := 0; i < 5; i++ {
		if err := sqlDB.Ping(); err != nil {
			fmt.Printf("⏳ Database connection attempt %d failed, retrying in 2s...\n", i+1)
			if i < 4 { // Don't sleep on last attempt
				time.Sleep(2 * time.Second)
			}
			continue
		}
		fmt.Println("✅ Database connection test successful")
		return nil
	}

	return fmt.Errorf("failed to connect to database after 5 attempts")
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
