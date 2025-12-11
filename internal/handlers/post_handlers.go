package handlers

import (
	"go-blog/internal/db"
	"go-blog/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePostRequest represents the request body for creating a post.
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}

// UpdatePostRequest represents the request body for updating a post.
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,min=1"`
}

// CreatePost handles creating a new blog post.
// Route: POST /api/posts
// Requirement: user must be authenticated (JWT), we use AuthMiddleware.
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	//Get current user ID from context(set by AuthMiddleware)
	userIdValue, exists := c.Get("userID")
	if !exists {
		c.JSON(500, gin.H{
			"error": "failed to get user from context",
		})
		return
	}

	userId, ok := userIdValue.(uint)
	if !ok {
		c.JSON(500, gin.H{
			"error": "invalid user ID type",
		})
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId,
	}

	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create post: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
		"userId":  post.UserID,
		"created": post.CreatedAt,
		"updated": post.UpdatedAt,
	})
}


// ListPosts handles fetching all posts.
// Route: GET /api/posts
// Requirement: public, no authentication needed.
func ListPosts(c *gin.Context) {
	var posts []models.Post

	// Preload user to include author info in response.
	if err := db.DB.Preload("User").Order("id DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to query posts",
		})
		return
	}

	// Build response to avoid exposing sensitive fields such as password.
	response := make([]gin.H, 0, len(posts))
	for _, p := range posts {
		response = append(response, gin.H{
			"id":      p.ID,
			"title":   p.Title,
			"content": p.Content,
			"userId":  p.UserID,
			"author": gin.H{
				"id":       p.User.ID,
				"username": p.User.Username,
				"email":    p.User.Email,
			},
			"created": p.CreatedAt,
			"updated": p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}


func GetPost(c *gin.Context) {
	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id parameter",
		})
		return
	}

	var post models.Post
	if err := db.DB.Preload("User").First(&post, uint(postID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "post not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
		"userId":  post.UserID,
		"author": gin.H{
			"id":       post.User.ID,
			"username": post.User.Username,
			"email":    post.User.Email,
		},
		"created": post.CreatedAt,
		"updated": post.UpdatedAt,
	})
}

// UpdatePost handles updating an existing post.
// Route: PUT /api/posts/:id
// Requirement: user must be authenticated, and must be the author of the post.
func UpdatePost(c *gin.Context) {
	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id parameter",
		})
		return
	}

	// Get current user ID from context.
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user in context",
		})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid userID type",
		})
		return
	}

	// Find the post.
	var post models.Post
	if err := db.DB.First(&post, uint(postID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "post not found",
		})
		return
	}

	// Check if current user is the author of the post.
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not the author of this post",
		})
		return
	}

	// Bind new title and content.
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := db.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
		"userId":  post.UserID,
		"created": post.CreatedAt,
		"updated": post.UpdatedAt,
	})
}

// DeletePost handles deleting an existing post.
// Route: DELETE /api/posts/:id
// Requirement: user must be authenticated, and must be the author of the post.
func DeletePost(c *gin.Context) {
	idParam := c.Param("id")
	postID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id parameter",
		})
		return
	}

	// Get current user ID from context.
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user in context",
		})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid userID type",
		})
		return
	}

	// Find the post.
	var post models.Post
	if err := db.DB.First(&post, uint(postID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "post not found",
		})
		return
	}

	// Check if current user is the author of the post.
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not the author of this post",
		})
		return
	}

	// Delete the post.
	if err := db.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": true,
		"id":      post.ID,
	})
}
