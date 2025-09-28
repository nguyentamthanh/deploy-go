package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("🧪 Testing database connection...")

	// Set test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "deploy_go")
	os.Setenv("DB_SSLMODE", "disable")

	// Test connection
	if err := ConnectDatabase(); err != nil {
		fmt.Printf("❌ Database connection test failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Database connection test successful!")
	CloseDatabase()
}
