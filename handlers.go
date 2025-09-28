package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User Handlers

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	var users []User

	if err := DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Users fetched successfully",
		Data:    users,
	})
}

// GetUser returns a specific user by ID
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "User fetched successfully",
		Data:    user,
	})
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	user := User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser updates an existing user
func UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch user",
			})
		}
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	if len(updates) > 0 {
		if err := DB.Model(&user).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to update user",
			})
			return
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser deletes a user (soft delete)
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch user",
			})
		}
		return
	}

	if err := DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

// Post Handlers

// GetPosts returns all posts with user information
func GetPosts(c *gin.Context) {
	var posts []Post

	if err := DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch posts",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Posts fetched successfully",
		Data:    posts,
	})
}

// GetPost returns a specific post by ID
func GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var post Post
	if err := DB.Preload("User").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "Post not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch post",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Post fetched successfully",
		Data:    post,
	})
}

// CreatePost creates a new post
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// Check if user exists
	var user User
	if err := DB.First(&user, req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to validate user",
			})
		}
		return
	}

	post := Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	}

	if err := DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create post",
		})
		return
	}

	// Load user information
	DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Post created successfully",
		Data:    post,
	})
}

// UpdatePost updates an existing post
func UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var post Post
	if err := DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "Post not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch post",
			})
		}
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	if len(updates) > 0 {
		if err := DB.Model(&post).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to update post",
			})
			return
		}
	}

	// Load user information
	DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Post updated successfully",
		Data:    post,
	})
}

// DeletePost deletes a post (soft delete)
func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var post Post
	if err := DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "Post not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch post",
			})
		}
		return
	}

	if err := DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

// GetUserPosts returns all posts by a specific user
func GetUserPosts(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Check if user exists
	var user User
	if err := DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch user",
			})
		}
		return
	}

	var posts []Post
	if err := DB.Preload("User").Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to fetch user posts",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "User posts fetched successfully",
		Data:    posts,
	})
}
