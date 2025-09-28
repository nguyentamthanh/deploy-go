package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Tạo Gin router
	r := gin.Default()

	// Route chính
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Xin chào! Đây là ứng dụng Golang cơ bản",
			"status":  "success",
		})
	})

	// Route health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Route API example
	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from API!",
			"version": "1.0.0",
		})
	})

	// Lắng nghe trên port 8080 (Render sẽ tự động set PORT env var)
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	r.Run(port)
}
