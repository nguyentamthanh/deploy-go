package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	fmt.Println("🚀 Starting application...")
	if err := ConnectDatabase(); err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		panic("Database connection failed: " + err.Error())
	}
	defer func() {
		fmt.Println("🔌 Closing database connection...")
		CloseDatabase()
	}()

	// Tạo Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Route chính
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Xin chào! Đây là ứng dụng Golang với PostgreSQL",
			"status":   "success",
			"database": "connected",
		})
	})

	// Route health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// Route API example
	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from API with PostgreSQL!",
			"version": "1.0.0",
		})
	})

	// User routes
	api := r.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("", GetUsers)          // GET /api/users
			users.GET("/:id", GetUser)       // GET /api/users/:id
			users.POST("", CreateUser)       // POST /api/users
			users.PUT("/:id", UpdateUser)    // PUT /api/users/:id
			users.DELETE("/:id", DeleteUser) // DELETE /api/users/:id
		}

		// Post routes
		posts := api.Group("/posts")
		{
			posts.GET("", GetPosts)                  // GET /api/posts
			posts.GET("/:id", GetPost)               // GET /api/posts/:id
			posts.POST("", CreatePost)               // POST /api/posts
			posts.PUT("/:id", UpdatePost)            // PUT /api/posts/:id
			posts.DELETE("/:id", DeletePost)         // DELETE /api/posts/:id
			posts.GET("/user/:userId", GetUserPosts) // GET /api/posts/user/:userId
		}
	}

	// Lắng nghe trên port 8080 (Render sẽ tự động set PORT env var)
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	r.Run(port)
}
