package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-blog/internal/db"
	"go-blog/internal/models"
)

// CreateCommentRequest represents the request body for creating a comment.
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"` // 评论内容不能为空
}

// CreateComment handles posting a comment to a specific post.
// Route: POST /api/posts/:id/comments
// Requirement: authentication needed.
func CreateComment(c *gin.Context) {
	// Parse post ID from URL path.
	postIDParam := c.Param("id")
	postID64, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid post ID",
		})
		return
	}
	postID := uint(postID64)

	// Check if the post exists.
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "post not found",
		})
		return
	}

	// Authentication: get userID from context.
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

	// Bind JSON request.
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  postID,
	}

	// Insert comment into DB.
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create comment",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      comment.ID,
		"content": comment.Content,
		"userId":  comment.UserID,
		"postId":  comment.PostID,
		"created": comment.CreatedAt,
	})
}

// ListComments handles fetching all comments for a specific post.
// Route: GET /api/posts/:id/comments
// Requirement: public access (no authentication needed).
func ListComments(c *gin.Context) {
	// Parse post ID.
	postIDParam := c.Param("id")
	postID64, err := strconv.ParseUint(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid post ID",
		})
		return
	}
	postID := uint(postID64)

	// Check if the post exists.
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "post not found",
		})
		return
	}

	// Query comments with user info.
	var comments []models.Comment
	if err := db.DB.Preload("User").Where("post_id = ?", postID).Order("id ASC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to query comments",
		})
		return
	}

	// Build response.
	response := make([]gin.H, 0, len(comments))
	for _, cm := range comments {
		response = append(response, gin.H{
			"id":      cm.ID,
			"content": cm.Content,
			"userId":  cm.UserID,
			"author": gin.H{
				"id":       cm.User.ID,
				"username": cm.User.Username,
				"email":    cm.User.Email,
			},
			"created": cm.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
